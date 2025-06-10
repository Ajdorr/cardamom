import { api } from '@core/api'
import { AxiosRequestConfig, AxiosResponse } from 'axios';
import { NavigateFunction } from 'react-router-dom';

const authorizedOAuthSources = ["github", "google", "facebook", "microsoft"]
var isOAuthLoggingIn = false

api.interceptors.request.use(onRequest, err => { return Promise.reject(err) })
api.interceptors.response.use(rsp => { return Promise.resolve(rsp) }, onResponseError)

export function login(email: string, password: string, callback: () => void) {
    api.post("auth/login", {
        email: email,
        password: password
    }).then(rsp => {
        onAuthenticate(rsp.headers["x-csrf-token"])
        callback()
    })
}

export function isAuthenticated(): boolean {
    let csrf = localStorage.getItem("csrf_token")
    if (!csrf) {
        return false
    }

    return true
}

export function startOAuthLogin(provider: string) {
    api.get(`/auth/oauth-start/${provider}`)
    .then(rsp => { window.location.assign(rsp.data.redirect_url) })
}

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

    api.post(`/auth/oauth-complete/${source}`, { code: code, state: state }).then(rsp => {

        onAuthenticate(rsp.headers["x-csrf-token"])
        const urlFwd = sessionStorage.getItem("auth_forward_pathname")
        if (urlFwd) {
            sessionStorage.removeItem("auth_forward_pathname")
            nav(urlFwd)
        } else {
            nav("/grocery")
        }
        isOAuthLoggingIn = false
    })

    return true
}

export function redirectIfNotAuthenticated(nav: NavigateFunction) {
    if (!isAuthenticated()) {
        refreshAuth().catch(e => {
            nav("/auth/login")
        })
    }
}

export enum LogLevel {
    ERROR = "Error",
    WARNING = "Warning"
}

export function log(message: string, logLevel: LogLevel, e: any) {
    api.post("/log", { level: logLevel, msg: message, data: e })
}

function onRequest(cfg: AxiosRequestConfig<any>) {
    // Add CSRF token to all requests
    if (cfg.headers) {
        let csrf = localStorage.getItem("csrf_token")
        if (csrf && csrf !== "") {
            cfg.headers["x-csrf-token"] = csrf
        }
    }
    return Promise.resolve(cfg)
}

function refreshAuth(): Promise<AxiosResponse<any, any>> {
    return new Promise<AxiosResponse<any, any>>((resolve, reject) => {
        api.post("auth/refresh").then(rsp => {
            let csrf = rsp.headers["x-csrf-token"]
            localStorage.setItem("csrf_token", csrf)
            resolve(rsp)
        }).catch(e => {
            if (!window.location.pathname.startsWith("/auth")) {
                sessionStorage.setItem("auth_forward_pathname", window.location.pathname)
            }
            reject(e)
        })
    })
}

function onResponseError(err: any) {
    if (err.config.url === "auth/refresh") {
        return Promise.reject(err)
    } else if (err.response.status === 401) {
        // Clear the CSRF token
        localStorage.removeItem("csrf_token")

        // If reauthentication was successful, retry the request
        return new Promise<AxiosResponse<any, any>>((resolve, reject) => {
            refreshAuth().then(rsp => {
                api.request(err.config).then(resolve).catch(reject)
            }).catch(e => {
                if (window.location.pathname !== "/auth/login") {
                    window.location.href = "/auth/login"
                }
                reject(err)
            })
        })
    } else {
        return Promise.reject(err)
    }
}

function onAuthenticate(csrf: string) {
    localStorage.setItem("csrf_token", csrf)
}