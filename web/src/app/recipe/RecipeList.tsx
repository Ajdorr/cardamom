import { useEffect, useState } from "react"
import { Link } from "react-router-dom"
import { api } from "../api"
import { ImageButton } from "../component/input"
import { RecipeModel } from "./schema"

function RecipeList() {
  const [recipes, setRecipes] = useState<RecipeModel[]>([])

  useEffect(() => {
    api.post("recipe/list").then(rsp => {
      setRecipes(rsp.data)
    }).catch(e => console.log(e))
  }, [])

  return (<div className="recipe-list-root">

    <div className="recipe-list-sub-menu theme-primary-light">
      <Link to="/recipe/create" >
        <img alt="New recipe" src="/icons/add-note.svg"
          id="recipe-index-create-btn"
          className="recipe-list-sub-menu-icon" />
      </Link>
      <Link to="/recipe/available" >
        <img alt="Get available recipes" src="/icons/shuffle.svg"
          id="recipe-index-get-available-btn"
          className="recipe-list-sub-menu-icon" />
      </Link>
      <Link to="/recipe/trash" >
        <img alt="view trashed recipes" src="/icons/delete.svg"
          id="recipe-index-get-trash-btn"
          className="recipe-list-sub-menu-icon" />
      </Link>
    </div>

    <div className="recipe-list-recipes">
      {recipes.length > 0 ? recipes.map(r => {
        return (
          <div key={r.uid} className="recipe-list-element-root">
            <Link to={`/recipe/${r.uid}`} >
              <div className="recipe-list-element-name"><span>{r.name}</span></div>
            </Link>
            <ImageButton alt="trash recipe" src="/icons/delete.svg" className="recipe-list-element-trash"
              onClick={e => {
                api.post("recipe/update", { uid: r.uid, is_trashed: true }).then(rsp => {
                  let newRecipes = [...recipes].filter(nr => { return nr.uid !== r.uid })
                  setRecipes(newRecipes)
                }).catch(console.log)
              }}
            />
          </div>
        )
      }) : <div className="recipe-list-empty">No recipes</div>}
    </div>
  </div>)
}

export default RecipeList