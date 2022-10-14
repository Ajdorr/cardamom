import "./recipe.css"
import { Component } from "react"
import { Link, Route, Routes } from "react-router-dom"
import RecipeList from "./RecipeList"
import RecipeSingle from "./RecipeSingle"
import RecipeShuffle from "./RecipeShuffle"
import RecipeTrashList from "./RecipeTrashList"
import RecipeSearch from "./RecipeSearch"

export function RecipeContextMenu() {
  return (
    <div className="recipe-list-sub-menu">
      <Link to="/recipe/search" >
        <img alt="Search for recipe" src="/icons/search.svg"
          id="recipe-index-search-btn" className="recipe-list-sub-menu-icon" />
      </Link>
      <Link to="/recipe/list" >
        <img alt="Recipe list" src="/icons/book.svg"
          id="recipe-index-list-btn" className="recipe-list-sub-menu-icon" />
      </Link>
      <Link to="/recipe/create" >
        <img alt="New recipe" src="/icons/add-note.svg"
          id="recipe-index-create-btn" className="recipe-list-sub-menu-icon" />
      </Link>
      <Link to="/recipe/available" >
        <img alt="Get available recipes" src="/icons/shuffle.svg"
          id="recipe-index-get-available-btn" className="recipe-list-sub-menu-icon" />
      </Link>
      <Link to="/recipe/trash" >
        <img alt="view trashed recipes" src="/icons/delete.svg"
          id="recipe-index-get-trash-btn" className="recipe-list-sub-menu-icon" />
      </Link>
    </div>
  )
}
class RecipeIndex extends Component {

  render() {
    return (<div className="recipe-index-root">

      <Routes>
        <Route path="/" element={<RecipeList />} />
        <Route path="list" element={<RecipeList />} />
        <Route path="create" element={<RecipeSingle isCreate={true} />} />
        <Route path="available" element={<RecipeShuffle />} />
        <Route path="trash" element={<RecipeTrashList />} />
        <Route path="search" element={<RecipeSearch />} />
        <Route path="edit/:recipeUid" element={<RecipeSingle isCreate={false} />} />
      </Routes>
    </div>)
  }

}

export default RecipeIndex