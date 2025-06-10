import "./auth.css"
import { useCallback, useState } from 'react';
import { useNavigate } from "react-router-dom"
import { api } from "../core/api";
import { login } from './authService';
import { FormText, FormPassword } from '../component/form'
import { ImageButton, TextButton } from '../component/input'
import { Theme } from '@core/theme';
import { startOAuthLogin } from '@auth/authService';

type LoginButtonProps = {
  id: string
  imgSrc: string
  provider: string
}

function LoginButton(props: LoginButtonProps) {

  return <div id={props.id} className="auth-login-btn-root"
    onClick={e => { startOAuthLogin(props.provider) }}>
    <img alt={`${props.provider} logo`} className="auth-login-btn-logo" src={props.imgSrc}></img>
    <div className="auth-login-btn-label format-font-small">
      <span>
        {`Sign in with ${props.provider.charAt(0).toUpperCase() + props.provider.substring(1).toLowerCase()}`}
      </span>
    </div>
  </div>
}

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
    <div className="auth-login-root theme-focus">
      <LoginButton id="auth-login-button-google" provider="google" imgSrc="/img/auth/google-logo.svg" />
      <LoginButton id="auth-login-button-microsoft" provider="microsoft" imgSrc="/img/auth/microsoft-logo.svg" />
      <LoginButton id="auth-login-button-facebook" provider="facebook" imgSrc="/img/auth/facebook-logo.png" />
      <LoginButton id="auth-login-button-github" provider="github" imgSrc="/img/auth/github-logo.svg" />

      <div id="auth-login-show-traditional">
        <ImageButton alt="show traditional login method" src="/icons/drop-down.svg"
          onClick={e => setShowTraditional(!showTraditional)} />
      </div>

      <div style={{ display: showTraditional ? "flex" : "none" }} className="auth-login-traditional-root">
        <FormText id="auth-login-traditional-email" label="Email" value={email} onChange={setEmail} />
        <FormPassword id="auth-login-traditional-password" label="Password" value={password} onChange={setPassword} />
        <TextButton id="auth-login-traditional-submit" label="Login" onClick={onSubmit} theme={Theme.Primary} />
      </div>
    </div>
  )
}

export default Login;