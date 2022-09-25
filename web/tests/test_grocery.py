from tests import login
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.remote.webelement import WebElement


def get_attr(e: WebElement, by: By, selector: str, attr: str) -> str:
  return e.find_element(by, selector).get_attribute(attr)


def test_grocery_basic_usage():
  d, w = login()

  # Go to grocery list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, "#workspace-menu-link-grocery img")).click()

  # Add items
  add_item = w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, "#grocery-list-add-item input"))
  add_item.send_keys("apples")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("peArs")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("banaNAs")

  # Add to Costco
  set_store = d.find_element(By.CSS_SELECTOR, "#grocery-list-store input")
  set_store.send_keys("Costco")
  add_item.send_keys("chiCken")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("eGGs")

  # Add to Metro
  clear_store = d.find_element(
      By.CSS_SELECTOR, "#grocery-list-store .modifiable-drop-down-clear img")
  clear_store.click()
  set_store.send_keys("Metro")
  add_item.send_keys("broCcolI")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("kale")

  # Show all
  clear_store.click()

  # verify
  groceries = dict(map(lambda g: (
      get_attr(g, By.CSS_SELECTOR, ".grocery-item-input input", "value"),
      get_attr(g, By.CSS_SELECTOR, ".grocery-item-store input", "value")),
      d.find_elements(By.CLASS_NAME, "grocery-item-root")))

  assert "apples" in groceries and groceries["apples"] == ""
  assert "pears" in groceries and groceries["pears"] == ""
  assert "bananas" in groceries and groceries["bananas"] == ""
  assert "chicken" in groceries and groceries["chicken"] == "costco"
  assert "eggs" in groceries and groceries["eggs"] == "costco"
  assert "broccoli" in groceries and groceries["broccoli"] == "metro"
  assert "kale" in groceries and groceries["kale"] == "metro"


def test_grocery_collect():
  d, w = login()

  # Go to grocery list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, "#workspace-menu-link-grocery img")).click()

  add_item: WebElement = w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, "#grocery-list-add-item input"))
  set_store = d.find_element(By.CSS_SELECTOR, "#grocery-list-store input")

  set_store.send_keys("Costco")
  add_item.send_keys("rice")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("flour")
  add_item.send_keys(Keys.ENTER)

  # find a grocery item element
  grocery_ele = "//div[contains(@class,'grocery-item-input')]/input[@value='%s']"
  # find associated collect button
  collect_btn = (
      grocery_ele + "/../../div[contains(@class,'grocery-item-collect')]/img")
  d.find_element(By.XPATH, collect_btn % "rice").click()

  # Verify collected
  collected_item = "//div[contains(@class, 'grocery-list-collected-root')]/span[text()='%s']"
  w.until(lambda x: x.find_element(By.XPATH, collected_item % "rice"))
  assert len(d.find_elements(By.XPATH, grocery_ele % "rice")) == 0
  assert d.find_element(
      By.XPATH, (collected_item % "rice") + "/../span[text()='costco']") != None

  # Uncollect
  d.find_element(
      By.XPATH, (collected_item % "rice") +
      "/../div/img").click()

  # Verify uncollect
  w.until(lambda x: x.find_element(By.XPATH, grocery_ele % "rice")) != None
  assert len(d.find_elements(By.XPATH, (collected_item % "rice"))) == 0

  # collect both
  d.find_element(By.XPATH, collect_btn % "rice").click()
  d.find_element(By.XPATH, collect_btn % "flour").click()

  # Verify collect
  w.until(lambda x: x.find_element(
      By.XPATH, (collected_item % "flour") + "/../span[text()='costco']"))
  assert d.find_element(
      By.XPATH, (collected_item % "rice") + "/../span[text()='costco']") != None
  assert len(d.find_elements(By.XPATH, grocery_ele % "flour")) == 0
  assert len(d.find_elements(By.XPATH, grocery_ele % "rice")) == 0

  # clear collected
  d.find_element(
      By.CSS_SELECTOR, ".grocery-list-collected-clear-all img").click()
  assert len(d.find_elements(By.XPATH, grocery_ele % "rice")) == 0
  assert len(d.find_elements(
      By.XPATH, (collected_item % "rice") + "/../span[text()='costco']")) == 0
  assert len(d.find_elements(By.XPATH, grocery_ele % "flour")) == 0
  assert len(d.find_elements(
      By.XPATH, (collected_item % "flour") + "/../span[text()='costco']")) == 0
