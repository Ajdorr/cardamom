import './format.css';
import './theme.css';
import './App.css';
import { Routes, Route, Link, useNavigate } from 'react-router-dom';
import Account from './auth/Account';
import { useEffect, useState } from 'react';
import { redirectIfNotAuthenticated } from './api';
import GroceryList from './grocery/GroceryList';
import InventoryList, { InventoryMenu as InventoryContextMenu } from './inventory/InventoryList';
import RecipeIndex, { RecipeContextMenu } from './recipe/RecipeIndex';
import AuthIndex from './auth/AuthIndex';
import { ImageButton } from './component/input';
import AppCache from './AppCache';

function Home() {
  return (
    <div className="home-root">
      <div className="home-greeting theme-primary">Welcome to Cardamom!</div>
      <Link to="/auth/login" className="home-login-link theme-primary">
        <div>Login or Sign Up</div>
      </Link>
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

function WorkspaceHeader() {

  const [showMenu, setShowMenu] = useState(false)
  const nav = useNavigate()

  return (
    <div className="workspace-menu-bar theme-primary" onMouseLeave={e => { if (showMenu) setShowMenu(false) }}>

      <ImageButton alt="Show menu" src="/icons/menu.svg" className="workspace-menu-bar-show-btn"
        onClick={e => { setShowMenu(!showMenu) }} />

      <div className="workspace-menu-bar-context-sensitive">
        <Routes>
          <Route path="recipe/:filter" element={<RecipeContextMenu />} />
          <Route path="recipe/*" element={<RecipeContextMenu />} />
          <Route path="inventory/*" element={<InventoryContextMenu />} />
          <Route path="inventory/:filter" element={<InventoryContextMenu />} />
        </Routes>

      </div>

      <div style={{ display: showMenu ? "flex" : "none" }} className="workspace-menu-bar-overlay theme-primary" >

        <ImageButton id="workspace-menu-link-grocery" alt="Go to grocery list" src="/icons/cart.svg"
          onClick={e => { setShowMenu(false); nav("/grocery") }} />
        <ImageButton id="workspace-menu-link-inventory" alt="Go to inventory list" src="/icons/inventory.svg"
          onClick={e => { setShowMenu(false); nav("/inventory") }} />
        <ImageButton id="workspace-menu-link-recipe" alt="Go to recipe list" src="/icons/book.svg"
          onClick={e => { setShowMenu(false); nav("/recipe/list") }} />
        <ImageButton id="workspace-menu-link-account" alt="Go to account settings" src="/icons/settings.svg"
          onClick={e => { setShowMenu(false); nav("/account") }} />

      </div>
    </div>
  )
}

function Workspace() {

  const nav = useNavigate()

  // eslint-disable-next-line
  useEffect(() => { redirectIfNotAuthenticated(nav) }, [])

  return (
    <div className="workspace-root">
      <WorkspaceHeader />
      <div className="workspace-main">
        <AppCache>
          <Routes>
            <Route path="grocery" element={<GroceryList />} />
            <Route path="inventory" element={<InventoryList />} />
            <Route path="inventory/:filter" element={<InventoryList />} />
            <Route path="recipe/*" element={<RecipeIndex />} />
            <Route path="account" element={<Account />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </AppCache>
      </div>
    </div>
  );
}

function App() {

  // Set title
  useEffect(() => { document.title = "Cardamom" }, [])

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
