import "./auth.css"
import { useState } from 'react';
import { useNavigate } from "react-router-dom"
import { login, startOAuthLogin } from '../api';
import { TextButton } from '../component/input'
import { FormText, FormPassword } from '../component/form'
import { Theme } from '../theme';


function Login() {
  // const [email, setEmail] = useState("")
  // const [password, setPassword] = useState("")
  // const nav = useNavigate()
  // const onSubmit = () => {
  //   login(email, password, () => {
  //     const fwd = sessionStorage.getItem("auth_forward_pathname")
  //     if (fwd) {
  //       sessionStorage.removeItem("auth_forward_pathname")
  //       nav(fwd)
  //     } else {
  //       nav("/grocery")
  //     }
  //   })
  // }

  return (
    <div className="auth-login-root theme-background">
      <div>
        <TextButton label="Login with Github" theme={Theme.None}
          onClick={e => { startOAuthLogin("github") }} />
      </div>
      <div>
        <TextButton label="Login with Google" theme={Theme.None}
          onClick={e => { startOAuthLogin("google") }} />
      </div>
      <div>
        <TextButton label="Login with Facebook" theme={Theme.None}
          onClick={e => { startOAuthLogin("facebook") }} />
      </div>
      <div>
        <TextButton label="Login with Microsoft" theme={Theme.None}
          onClick={e => { startOAuthLogin("microsoft") }} />
      </div>

      {/* <div>
        <FormText label="Email" value={email} onChange={setEmail} />
        <FormPassword label="Password" value={password} onChange={setPassword} />
        <TextButton label="Login" onClick={onSubmit} theme={Theme.Primary} />
      </div> */}
    </div>
  )
}

export default Login;