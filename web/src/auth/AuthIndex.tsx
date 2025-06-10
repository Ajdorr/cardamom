import { useEffect, useState } from "react";
import { Route, Routes, useNavigate, useParams, useSearchParams } from "react-router-dom";
import { completeOAuthLogin } from "./authService";
import Login from "./Login";

function OAuthReturn() {

  const nav = useNavigate()
  let { source } = useParams()
  const [searchParams,] = useSearchParams()
  const [displayText, setDisplayText] = useState("Redirecting...")
  const code = searchParams.get("code")
  const state = searchParams.get("state")

  useEffect(() => {
    if (!(source && code && state && completeOAuthLogin(nav, source, code, state))) {
      setDisplayText("Something went wrong, please try again later")
    } // eslint-disable-next-line
  }, [source, code, state, completeOAuthLogin, setDisplayText])

  return <div className="auth-oauth-return-root">{displayText}</div>
}

export default function AuthIndex() {
  return (
    <div className="auth-index-root theme-background">
      <div className="auth-index-workspace theme-focus">
        <div className="auth-index-branding">
          <img alt="logo" src="/img/logo/logo_300.png" className="auth-index-branding-img" />
          <div className="auth-index-branding-title format-font-header-medium">Cardamom</div>
          <div className="auth-index-branding-tagline format-font-small">Grocery lists and recipe books, redesigned</div>
        </div>
        <Routes>
          <Route path="login" element={<Login />} />
          <Route path="oauth-return/:source" element={<OAuthReturn />} />
        </Routes>
      </div>
    </div>
  )
}