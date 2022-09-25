import "./auth.css"
import { useState } from 'react';
import { useNavigate } from "react-router-dom"
// import { startOAuthLogin } from '../api';
import { login, startOAuthLogin } from '../api';
import { FormText, FormPassword } from '../component/form'
import { ImageButton, TextButton } from '../component/input'
import { Theme } from '../theme';


function Login() {
  const [showTraditional, setShowTraditional] = useState(false)
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const nav = useNavigate()

  const onSubmit = () => {
    login(email, password, () => {
      const fwd = sessionStorage.getItem("auth_forward_pathname")
      if (fwd) {
        sessionStorage.removeItem("auth_forward_pathname")
        nav(fwd)
      } else {
        nav("/grocery")
      }
    })
  }

  return (
    <div className="auth-login-root theme-background">
      <div>
        <TextButton label="Login with Github" id="auth-login-button-github" theme={Theme.None}
          onClick={e => { startOAuthLogin("github") }} />
      </div>
      <div>
        <TextButton label="Login with Google" id="auth-login-button-google"
          theme={Theme.None} onClick={e => { startOAuthLogin("google") }} />
      </div>
      <div>
        <TextButton label="Login with Facebook" id="auth-login-button-facebook"
          theme={Theme.None} onClick={e => { startOAuthLogin("facebook") }} />
      </div>
      <div>
        <TextButton label="Login with Microsoft" id="auth-login-button-microsoft"
          theme={Theme.None} onClick={e => { startOAuthLogin("microsoft") }} />
      </div>

      <div id="auth-login-show-traditional">
        <ImageButton alt="show traditional login method" src="/icons/drop-down.svg"
          onClick={e => setShowTraditional(!showTraditional)} />
      </div>

      <div style={{ display: showTraditional ? "block" : "none" }}>
        <FormText id="auth-login-traditional-email" label="Email" value={email} onChange={setEmail} />
        <FormPassword id="auth-login-traditional-password" label="Password" value={password} onChange={setPassword} />
        <TextButton id="auth-login-traditional-submit" label="Login" onClick={onSubmit} theme={Theme.Primary} />
      </div>
    </div>
  )
}

export default Login;