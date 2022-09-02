import axios from 'axios';
import { NavigateFunction } from 'react-router-dom';

export const api = axios.create({
  baseURL: '/api',
  timeout: 60000,
})

api.interceptors.request.use(cfg => {
  if (cfg.headers) {
    var csrf = localStorage.getItem("csrf_token")
    if (csrf && csrf !== "") {
      cfg.headers["x-csrf-token"] = csrf
    }
  }
  return Promise.resolve(cfg)
}, err => Promise.reject(err))

api.interceptors.response.use(cfg => Promise.resolve(cfg),
  err => {
    if (err.response.status === 401) {
      localStorage.removeItem("csrf_token")
      if (!window.location.pathname.startsWith("/auth")) {
        sessionStorage.setItem("auth_forward_pathname", window.location.pathname)
      }
      window.location.href = "/auth/login"
    }
    return Promise.reject(err)
  })

export function login(email: string, password: string, callback: () => void) {
  api.post("auth/login", {
    email: email,
    password: password
  }).then(rsp => {
    onAuthenticate(rsp.headers["x-csrf-token"])
    callback()
  }).catch(e => {
    console.log(e) // FIXME need a logging system
  })
}

export function logout() {
  api.post("/auth/logout").then(rsp => {
    localStorage.removeItem("csrf_token")
    stopRefreshLoop()
    document.location.href = '/auth/login'
  })
}

function onAuthenticate(csrf: string) {
  localStorage.setItem("csrf_token", csrf)
  startRefreshLoop()
}

export function isAuthenticated(): boolean {
  var csrf = localStorage.getItem("csrf_token")
  if (!csrf) {
    return false
  }

  return true
}

export function startOAuthLogin(provider: string, forward_url?: string) {
  if (forward_url) {
    sessionStorage.setItem("auth_forward_pathname", forward_url)
  }

  api.get(`/auth/oauth2-start/${provider}`).then(rsp => {
    window.location.href = rsp.data.redirect_url
  }).catch(console.log)
}

const authorizedOAuthSources = ["github", "google", "facebook", "microsoft"]
var isOAuthLoggingIn = false

export function completeOAuthLogin(
  nav: NavigateFunction,
  source: string, code: string, state: string): boolean {

  if (source && authorizedOAuthSources.indexOf(source) < 0) {
    return false
  }

  if (isOAuthLoggingIn) {
    return true
  }
  isOAuthLoggingIn = true

  api.post(`/auth/oauth2-complete/${source}`, { code: code, state: state }).then(rsp => {

    onAuthenticate(rsp.headers["x-csrf-token"])
    const urlFwd = sessionStorage.getItem("auth_forward_pathname")
    if (urlFwd) {
      sessionStorage.removeItem("auth_forward_pathname")
      nav(urlFwd)
    } else {
      nav("/grocery")
    }
    isOAuthLoggingIn = false
  }).catch(e => console.log)
  return true
}

export function redirectIfNotAuthenticated(nav: NavigateFunction) {

  if (!isAuthenticated()) {
    sessionStorage.setItem("auth_forward_pathname", window.location.pathname)
    nav("/auth/login")
  }
}

// const refreshLoopMs = 3 * 1000
const refreshLoopMs = 55 * 60 * 1000
var refreshLoop: number

function startRefreshLoop() {
  stopRefreshLoop()
  refreshLoop = window.setInterval(refreshAuth, refreshLoopMs)
}

function stopRefreshLoop() {
  window.clearInterval(refreshLoop)
}

function refreshAuth() {
  api.post("auth/refresh").then(rsp => {
    var csrf = rsp.headers["x-csrf-token"]
    localStorage.setItem("csrf_token", csrf)
  }).catch(e => {
    stopRefreshLoop()
    console.log(e)
  })
}