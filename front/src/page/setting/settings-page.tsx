import { useTranslation } from 'react-i18next'

export const SettingsPage = () => {
  const { t } = useTranslation()
  return <div>{t('page.settings')}</div>
}

/* ========== 설정 페이지 TODO ==========
 *
 * [1] 멤버 초대
 *   - [Backend] 그룹 기준 발급된 초대 토큰 목록 API (ListInviteTokens)
 *   - [Backend] 초대 토큰 취소/무효화 API (RevokeInviteToken) — 선택
 *   - [Front] 초대 링크 생성 UI (클랜 선택 → 생성 → URL 복사) — POST /invite 연동됨
 *   - [Front] 발급한 초대 목록 표시 (만료일, 취소 버튼 등)
 *
 * [2] 그룹명 변경
 *   - [Backend] UpdateGroup (groups.name) — query, store, service, handler
 *   - [Front] 그룹 설정 섹션: 그룹명 입력 폼 + 저장
 *
 * [3] 멤버 이름(호칭) 변경
 *   - [Backend] 그룹/클랜 기준 멤버 목록 API (ListMembersByGroup 또는 ByClan)
 *   - [Backend] 사용자 표시명(users.name) 수정 API (UpdateUserName) — 본인 또는 관리자
 *   - [Front] 멤버 목록 + 이름 편집 UI
 *
 * [4] 멤버 소속 클랜 변경
 *   - [Backend] 사용자 clan_group_id 수정 API (UpdateUserClanGroup) — 권한: 동일 그룹 내 관리자
 *   - [Front] 멤버 목록에서 클랜 선택 드롭다운
 *
 * [5] 멤버 권한 설정
 *   - [Backend] 권한 모델 정리: clan_groups.is_admin vs 사용자별 역할(필요 시 users에 role 추가)
 *   - [Backend] 관리자/일반 전환 API (UpdateClanGroup is_admin 또는 UpdateUserRole)
 *   - [Front] 멤버 목록에서 권한(관리자/일반) 토글 또는 선택
 *
 * [6] 공통/인프라
 *   - [Backend] 설정 관련 API에 권한 가드: 그룹 관리자만 그룹명/멤버/초대 관리 가능
 *   - [Backend] 본인 정보 수정(이름)은 본인만 가능하도록 가드
 *   - [Front] 설정 페이지 레이아웃: 탭 또는 섹션(그룹 정보 / 멤버·초대 / 클랜 등)
 *   - [Front] i18n: 설정 페이지 라벨/메시지 추가
 *
 * [7] 선택 사항
 *   - 그룹 삭제 (데이터 정책: 소프트 삭제, 멤버/앨범 처리)
 *   - 클랜 추가/삭제/이름 변경 (CreateClanGroup 이미 있음, UpdateClanGroup 추가)
 *   - 그룹 나가기/탈퇴 (본인만, 마지막 멤버 시 방지 등)
 */