// 캐시 이름 제공 및 모든 페이지 즉시 제어를 위한 모듈
import { cacheNames, clientsClaim } from 'workbox-core'
// 라우트 등록 및 기본/에러 핸들러 설정을 위한 라우팅 모듈
import { setCatchHandler, setDefaultHandler } from 'workbox-routing'
// 캐싱 전략을 위한 캐싱 모듈
import { NetworkFirst, NetworkOnly } from 'workbox-strategies'

// Firebase (주석 처리됨)
// declare let firebase: any;
// importScripts('https://www.gstatic.com/firebasejs/9.6.8/firebase-app-compat.js');
// importScripts('https://www.gstatic.com/firebasejs/9.6.8/firebase-messaging-compat.js');
// import { initializeApp } from 'firebase/app'
// import { getMessaging, isSupported } from 'firebase/messaging/sw'

// __WB_MANIFEST 속성을 포함하도록 ServiceWorkerGlobalScope 확장
interface MyServiceWorkerGlobalScope extends ServiceWorkerGlobalScope {
  __WB_MANIFEST: any
}

// TypeScript에 올바른 전역 타입 제공
declare let self: MyServiceWorkerGlobalScope

// install 및 activate 이벤트에서 사용할 ExtendableEvent 타입 선언
declare type ExtendableEvent = any

const data = {
  race: false, // 네트워크 우선, 실패 시 캐시된 응답 반환
  debug: false, // 가로채진 요청/응답에 대한 디버그 메시지 로깅 비활성화
  credentials: 'same-origin', // 같은 출처의 리소스만 요청
  networkTimeoutSeconds: 0, // 네트워크 요청 타임아웃(초), 0은 타임아웃 없음
  fallback: 'index.html', // 요청 실패 시 index.html로 폴백
}

// 앱에서 사용할 미리 정의된 캐시 이름
const precacheCacheName = cacheNames.precache
const runtimeCacheName = cacheNames.runtime

// 매니페스트 가져오기
// workbox가 서비스 워커에 매니페스트를 주입하기 위해 필요함
// __WB_MANIFEST는 vite-plugin-pwa가 빌드 시 주입함
const manifest = self.__WB_MANIFEST || []
const cacheEntries: (string | { url: string; revision?: string | null })[] = manifest

// 개발 모드 확인 (개발 모드에서는 매니페스트가 비어있음)
const isDevMode = cacheEntries.length === 0

// 서비스 워커가 처음 설치될 때 리소스 캐싱
self.addEventListener('install', (event: ExtendableEvent) => {
  // 개발 모드에서는 캐싱 건너뛰기
  if (isDevMode) {
    data.debug && console.log('개발 모드: 캐시 설치 건너뛰기')
    return
  }

  // 브라우저는 Promise가 해결될 때까지 대기
  event.waitUntil(
    // precache를 열고 매니페스트의 모든 리소스를 캐싱
    caches.open(precacheCacheName).then((cache) => {
      // 매니페스트 항목에서 URL 추출 (문자열 또는 url 속성을 가진 객체)
      const urlsToCache = cacheEntries.map((entry) => {
        return typeof entry === 'string' ? entry : entry.url
      })
      return cache.addAll(urlsToCache)
    }),
  )
})

// 서비스 워커 활성화 시, 매니페스트에 없는 오래된 캐시 정리
self.addEventListener('activate', (event: ExtendableEvent) => {
  // 개발 모드에서는 캐시 정리 건너뛰기
  if (isDevMode) {
    data.debug && console.log('개발 모드: 캐시 정리 건너뛰기')
    // 개발 모드에서는 모든 캐시를 삭제하여 최신 데이터 보장
    event.waitUntil(
      caches.keys().then((cacheNamesList) => {
        return Promise.all(
          cacheNamesList.map((cacheName) => {
            data.debug && console.log(`개발 모드: 캐시 삭제: ${cacheName}`)
            return caches.delete(cacheName)
          }),
        )
      }),
    )
    return
  }

  // 오래된 캐시 정리
  event.waitUntil(
    Promise.all([
      // 현재 캐시 이름에 없는 모든 캐시 삭제
      caches.keys().then((cacheNamesList) => {
        const validCacheNames = [precacheCacheName, runtimeCacheName]
        return Promise.all(
          cacheNamesList.map((cacheName) => {
            if (!validCacheNames.includes(cacheName)) {
              data.debug && console.log(`오래된 캐시 삭제: ${cacheName}`)
              return caches.delete(cacheName)
            }
            return Promise.resolve(false)
          }),
        )
      }),
      // precache 캐시에서 오래된 항목 정리
      caches.open(precacheCacheName).then((cache) => {
        const currentUrls = new Set(
          cacheEntries.map((entry) => {
            return typeof entry === 'string' ? entry : entry.url
          }),
        )
        return cache.keys().then((keys) => {
          return Promise.all(
            keys.map((request) => {
              const url = request.url
              // 매니페스트에 더 이상 없는 캐시 항목 제거
              if (!currentUrls.has(url)) {
                data.debug && console.log(`오래된 precache 항목 제거: ${url}`)
                return cache.delete(request)
              }
              return Promise.resolve(false)
            }),
          )
        })
      }),
    ]),
  )
})

// 인식되지 않은 모든 라우트를 네트워크에서 가져오도록 설정
// 개발 모드에서는 NetworkOnly 전략을 사용하여 캐시를 완전히 우회
if (isDevMode) {
  // 개발 모드에서는 항상 네트워크에서만 가져오기 (캐시 사용 안 함)
  setDefaultHandler(new NetworkOnly())
} else {
  setDefaultHandler(new NetworkFirst())
}

// 네트워크에서 리소스를 가져올 수 없을 때 호출되는 메서드
setCatchHandler(({ event }: any): Promise<Response> => {
  // 개발 모드에서는 캐시 fallback 사용 안 함
  if (isDevMode) {
    return Promise.resolve(Response.error())
  }

  switch (event.request.destination) {
    case 'document':
      return caches.match(data.fallback).then((r) => {
        return r ? Promise.resolve(r) : Promise.resolve(Response.error())
      })
    default:
      return Promise.resolve(Response.error())
  }
})

// 새 서비스 워커가 skipWaiting 상태로 유지되므로 필요함
// 활성화되지 않으면 캐시가 정리되지 않음
self.skipWaiting()
clientsClaim()

// const config = {
//   apiKey: import.meta.env['VITE_API_KEY'],
//   authDomain: import.meta.env['VITE_AUTH_DOMAIN'],
//   projectId: import.meta.env['VITE_PROJECT_ID'],
//   storageBucket: import.meta.env['VITE_STORAGE_BUCKET'],
//   messagingSenderId: import.meta.env['VITE_MESSAGING_SENDER_ID'],
//   appId: import.meta.env['VITE_APP_ID'],
// }

// const app = initializeApp(config)

// let messages: string[] = []

// const braodcast = new BroadcastChannel('ch-notice')

// braodcast.onmessage = (event) => {
//   if (event.data.type === 'getMessages') {
//     braodcast.postMessage(messages)
//     messages = []
//   }
// }

// isSupported().then((supported) => {
//   if (supported) {
//     getMessaging(app)
//   }
// })
