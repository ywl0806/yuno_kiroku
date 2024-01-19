//Import Cache Names and ClientClaim module for providing the cache name and taking control of all pages immediately
import { cacheNames, clientsClaim } from 'workbox-core'

//Import routing modules for registering routes and setting default and catch handlers
import { setCatchHandler, setDefaultHandler } from 'workbox-routing'

//Import caching modules for caching strategies
import { NetworkFirst } from 'workbox-strategies'

//Import module for caching precached assets

//Firebase
// declare let firebase: any;
// importScripts('https://www.gstatic.com/firebasejs/9.6.8/firebase-app-compat.js');
import { initializeApp } from 'firebase/app'
// importScripts('https://www.gstatic.com/firebasejs/9.6.8/firebase-messaging-compat.js');
import { getMessaging, isSupported } from 'firebase/messaging/sw'

//Extend the ServiceWorkerGlobalScope to include the __WB_MANIFEST property
interface MyServiceWorkerGlobalScope extends ServiceWorkerGlobalScope {
  __WB_MANIFEST: any
}

// Give TypeScript the correct global.
declare let self: MyServiceWorkerGlobalScope

// Declare type for ExtendableEvent to use in install and activate events
declare type ExtendableEvent = any

const data = {
  race: false, //Fetch first, if it fails, return a previously cached response
  debug: false, //Don't log debug messages for intercepted requests and responses
  credentials: 'same-origin', //Only request resources from the same origin
  networkTimeoutSeconds: 0, //Timout in seconds for network requests; 0 means no timeout
  fallback: 'index.html', //Fallback to index.html if the request fails
}

// Access the pre-defined cache names for use in this app
const cacheName = cacheNames.runtime

//Retrieve the manifest. First define asynchronus function to retrieve the manifest
// This is also required for the injection of the manifest into the service worker by workbox
// So despite it being outdate, Your application will not build without it
//Array for resources that have been cached by the service worker
const cacheEntries: RequestInfo[] = []

// Cache resources when the service worker is first installed
self.addEventListener('install', (event: ExtendableEvent) => {
  // The browser will wait until the promise is resolved
  event.waitUntil(
    // Open the cache and cache all the resources in it. This may include resources
    // that are not in the manifest
    caches.open(cacheName).then((cache) => {
      return cache.addAll(cacheEntries)
    }),
  )
})

// Upon activating the service worker, clear outdated caches by removing caches associated with
// URL resources that are not in the manifest URL array
self.addEventListener('activate', (event: ExtendableEvent) => {
  // - clean up outdated runtime cache
  event.waitUntil(
    caches.open(cacheName).then((cache) => {
      // clean up those who are not listed in manifestURLs since the manifestURLs are the only
      // resources that are unlikely to be outdated
      cache.keys().then((keys) => {
        keys.forEach((request) => {
          data.debug && console.log(`Checking cache entry to be removed: ${request.url}`)
        })
      })
    }),
  )
})

// Inform the service worker to send all routes that are not recognized to the network to be fetched
setDefaultHandler(new NetworkFirst())

// This method is called when the service worker is unable to fetch a resource from the network
setCatchHandler(({ event }: any): Promise<Response> => {
  switch (event.request.destination) {
    case 'document':
      return caches.match(data.fallback).then((r) => {
        return r ? Promise.resolve(r) : Promise.resolve(Response.error())
      })
    default:
      return Promise.resolve(Response.error())
  }
})

// this is necessary, since the new service worker will keep on skipWaiting state
// and then, caches will not be cleared since it is not activated
self.skipWaiting()
clientsClaim()

const config = {
  apiKey: import.meta.env['VITE_API_KEY'],
  authDomain: import.meta.env['VITE_AUTH_DOMAIN'],
  projectId: import.meta.env['VITE_PROJECT_ID'],
  storageBucket: import.meta.env['VITE_STORAGE_BUCKET'],
  messagingSenderId: import.meta.env['VITE_MESSAGING_SENDER_ID'],
  appId: import.meta.env['VITE_APP_ID'],
}

const app = initializeApp(config)

// let messages: string[] = []

// const braodcast = new BroadcastChannel('ch-notice')

// braodcast.onmessage = (event) => {
//   if (event.data.type === 'getMessages') {
//     braodcast.postMessage(messages)
//     messages = []
//   }
// }

isSupported().then((supported) => {
  if (supported) {
    getMessaging(app)
  }
})
