import './format.css';
import './theme.css';
import './App.css';
import { Routes, Route, Link, useNavigate } from 'react-router-dom';
import Account from './auth/Account';
import { useEffect } from 'react';
import { redirectIfNotAuthenticated } from './api';
import GroceryList from './grocery/GroceryList';
import InventoryList from './inventory/InventoryList';
import RecipeIndex from './recipe/RecipeIndex';
import AuthIndex from './auth/AuthIndex';

function Home() {
  return (
    <div className="home-root">
      <div>Welcome to Cardamom!</div>
      <Link to="/auth/login"><div>Login</div></Link>
    </div>
  );
}

function NotFound() {
  return (
    <div className="not-found-root">
      <div>404 Page not found</div>
      <div>You will only find anthropomorphic dragons here...</div>
    </div>
  );
}

function Workspace() {

  const nav = useNavigate()

  // eslint-disable-next-line
  useEffect(() => { redirectIfNotAuthenticated(nav) }, [])

  return (
    <div className="workspace-root">
      <div className="workspace-menu-bar theme-primary">
        <Link to="/grocery"><img src="/icons/cart.svg" alt="grocery"/></Link>
        <Link to="/inventory"><img src="/icons/inventory.svg" alt="inventory" /></Link>
        <Link to="/recipe"><img src="/icons/book.svg" alt="recipes" /></Link>
        <Link to="/account"><img src="/icons/menu.svg" alt="grocery" /></Link>
        {/* Account */}
        {/* <TextButton label="Logout" theme={Theme.Primary} onClick={e => logout()} /> */}
      </div>
      <div className="workspace-main">
        <Routes>
          <Route path="grocery" element={<GroceryList />} />
          <Route path="inventory" element={<InventoryList />} />
          <Route path="recipe/*" element={<RecipeIndex />} />
          <Route path="account" element={<Account />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </div>
    </div>
  );
}

function App() {
  return (<div className="app-root">
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/auth/*" element={<AuthIndex />} />
      <Route path="/*" element={<Workspace />} />
    </Routes>
  </div>
  )
}

export default App;
