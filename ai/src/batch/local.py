"""
로컬 개발용 진입점

face_recognition_jobs 테이블을 폴링하여 job을 가져옵니다.
SQS 없이 DB만으로 동작하며, IDLE_EXIT_SECONDS 동안 job이 없으면 종료합니다.
"""

import logging
import os
import signal
import time

import psycopg2
import psycopg2.extras

from src.batch.face_recognition_batch import get_s3_client, logger, process_job

DATABASE_URL = os.environ.get("DATABASE_URL", "")

POLL_INTERVAL = 5
BATCH_SIZE = 10
IDLE_EXIT_SECONDS = 60

_running = True


def get_db_conn():
    return psycopg2.connect(DATABASE_URL)


def fetch_pending_jobs(conn, limit: int) -> list:
    """pending job을 원자적으로 fetch (FOR UPDATE SKIP LOCKED)"""
    with conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as cur:
        cur.execute(
            """
            UPDATE face_recognition_jobs
            SET status = '02',
                attempt_count = attempt_count + 1,
                updated_at = NOW()
            WHERE id IN (
                SELECT id FROM face_recognition_jobs
                WHERE status = '01'
                ORDER BY created_at ASC
                LIMIT %s
                FOR UPDATE SKIP LOCKED
            )
            RETURNING *
            """,
            (limit,),
        )
        jobs = cur.fetchall()
    conn.commit()
    return [dict(j) for j in jobs]


def run_local_batch_loop():
    logger.info("AI Batch Worker (local/DB) 시작")
    conn = get_db_conn()
    s3_client = get_s3_client()
    idle_start: float | None = None

    while _running:
        try:
            jobs = fetch_pending_jobs(conn, BATCH_SIZE)
            if jobs:
                idle_start = None
                logger.info(f"처리할 job: {len(jobs)}개")
                for job in jobs:
                    if not _running:
                        break
                    process_job(s3_client, job)
            else:
                now = time.monotonic()
                if idle_start is None:
                    idle_start = now
                elif now - idle_start >= IDLE_EXIT_SECONDS:
                    logger.info(f"{IDLE_EXIT_SECONDS}초 동안 pending job 없음 - 종료")
                    break
                time.sleep(POLL_INTERVAL)
        except psycopg2.OperationalError:
            logger.warning("DB 연결 재시도...")
            try:
                conn.close()
            except Exception:
                pass
            time.sleep(5)
            conn = get_db_conn()
        except Exception as e:
            logger.error(f"배치 루프 에러: {e}")
            time.sleep(POLL_INTERVAL)

    logger.info("AI Batch Worker (local/DB) 종료")
    conn.close()


def _handle_signal(sig, frame):
    global _running
    logger.info("종료 신호 수신 - 현재 job 완료 후 종료")
    _running = False


if __name__ == "__main__":
    signal.signal(signal.SIGTERM, _handle_signal)
    signal.signal(signal.SIGINT, _handle_signal)
    run_local_batch_loop()
