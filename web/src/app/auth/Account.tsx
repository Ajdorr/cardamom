import { logout } from "../api"
import { TextButton } from "../component/input"
import { Theme } from "../theme"

function Account() {

  return (<div className="auth-account-root theme-surface">
    <div className="auth-account-logout theme-focus">
      <TextButton label="Logout" theme={Theme.Primary} onClick={e => logout()} />
    </div>
  </div>)
}

export default Account