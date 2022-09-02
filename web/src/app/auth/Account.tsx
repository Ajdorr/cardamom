import { useState } from "react"
import { api, logout } from "../api"
import { FormPassword } from "../component/form"
import { TextButton } from "../component/input"
import { Theme } from "../theme"

function Account() {

  const [currentPassword, setCurrentPassword] = useState("")
  const [newPassword, setNewPassword] = useState("")
  const [verifyPassword, setVerifyPassword] = useState("")
  const [errMsg, setErrMsg] = useState("")

  return (<div className="auth-account-root theme-surface">
    <div className="auth-account-logout theme-focus">
      <TextButton label="Logout" theme={Theme.Primary} onClick={e => logout()} />
    </div>
    <div className="auth-account-change-password theme-focus">
      <div className="auth-account-change-password-error">{errMsg}</div>
      <div className="auth-account-change-password-error format-text-small">{errMsg}</div>
      <FormPassword value={currentPassword} label="Current Password" onChange={setCurrentPassword} />
      <FormPassword value={newPassword} label="New Password" onChange={setNewPassword} />
      <FormPassword value={verifyPassword} label="Verify Password" onChange={setVerifyPassword} />
      <TextButton label="Change Password" theme={Theme.Primary}
        onClick={e => {
          if (newPassword !== verifyPassword) {
            setErrMsg("Passwords must match")
            return
          }

          api.post("/auth/set-password", {
            current_password: currentPassword,
            new_password: newPassword
          }).then(rsp => {
            setErrMsg("")
          }).catch(err => {
            setErrMsg("Unable to set password")
            console.log(e)
          })

        }} />
    </div>
  </div>)
}

export default Account