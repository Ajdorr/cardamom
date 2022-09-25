import { useEffect } from "react";
import { Route, Routes, useNavigate, useParams, useSearchParams } from "react-router-dom";
import { completeOAuthLogin } from "../api";
import Login from "./Login";

function OAuthReturn() {

  const nav = useNavigate()
  var { source } = useParams()
  const [searchParams,] = useSearchParams()
  const code = searchParams.get("code")
  const state = searchParams.get("state")

  useEffect(() => {
    if (source && code && state) {
      if (completeOAuthLogin(nav, source, code, state)) {
        // FIXME
      } else {
        // FIXME
      }
    }
  })

  return <div>Redirecting...</div>
}

function AuthIndex() {
  return (
    <div className="auth-index-root">
      <Routes>
        <Route path="login" element={<Login />} />
        <Route path="oauth-return/:source" element={<OAuthReturn />} />
      </Routes>
    </div>
  )
}

export default AuthIndex