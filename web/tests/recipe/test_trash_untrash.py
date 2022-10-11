import os
import json
from typing import Any
from venv import create
from tests import login
from . import verify
from .test_create import create_flow
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys

info = {
    "name": "Fried Rice",
    "desc": "Rice any time",
    "meal": "Dinner",
    "ingre": [
        {
            "quantity": "2",
            "unit": "cup",
            "item": "rice"
        }
    ],
    "instr": [
        "Grease wok",
        "fry rice"
    ]
}


def test_trash_untrash():
  d, w = login()

  # Go to recipe list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-recipe img").click()

  # Wait until screen loads
  w.until(lambda x: x.find_elements(
      By.CSS_SELECTOR, ".recipe-list-element-name"))

  recipe_eles = d.find_elements(
      By.XPATH,
      f"//div[contains(@class, 'recipe-list-element-name')]/span[text()='{info['name']}']")
  if len(recipe_eles) == 0:
    create_flow(d, w, info)

    # Go to recipe list screen
    w.until(lambda x: x.find_element(
        By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
    d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-recipe img").click()

    # Wait until screen loads
    w.until(lambda x: x.find_elements(
        By.CSS_SELECTOR, ".recipe-list-element-name"))

  d.find_element(
      By.XPATH,
      f"//div[contains(@class, 'recipe-list-element-name')]/span[text()='{info['name']}']"
      "/../../../div[contains(@class, 'recipe-list-element-trash')]/img"
  ).click()

  # Go to trash
  d.find_element(By.ID, "recipe-index-get-trash-btn").click()

  # Wait for list to populate
  w.until(lambda x: x.find_elements(
      By.CLASS_NAME, "recipe-trash-list-recipes"))

  d.find_element(
      By.XPATH,
      f"//div[contains(@class, 'recipe-trash-list-element-name')]/span[text()='{info['name']}']"
      "/../../div[contains(@class, 'recipe-list-element-untrash')]/img"
  ).click()

  # Go to recipe list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-recipe img").click()

  w.until(lambda x: x.find_element(
      By.XPATH,
      f"//div[contains(@class, 'recipe-list-element-name')]/span[text()='{info['name']}']"
  ))
