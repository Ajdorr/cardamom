from tests import login, clear
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.remote.webelement import WebElement
from selenium.webdriver.support.select import Select

inventory_items = [
    "potato", "taco", "rice", "chicken", "cumin"
]

inv_item_xpath = "//div[contains(@class,'inventory-item-input')]/input[@value='%s']"


def test_required_data():
  d, w = login()

  # Go to inventory list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-inventory img").click()

  add_item = d.find_element(
      By.CSS_SELECTOR, ".inventory-list-add-item-input input")

  for inv in inventory_items:
    add_item.send_keys(inv)
    add_item.send_keys(Keys.ENTER)

  for inv in inventory_items:
    eles = d.find_elements(
        By.XPATH,
        f"//div[contains(@class,'inventory-item-input')]/input[@value='{inv}']")
    assert len(eles) == 1


def test_basic_usage():
  d, w = login()

  # Go to inventory list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-inventory img").click()

  # Cleanup previous runs
  del_eles = d.find_elements(By.XPATH, inv_item_xpath % "pasta") + \
      d.find_elements(By.XPATH, inv_item_xpath % "strawberries")
  for ele in del_eles:
    ele.find_element(
        By.XPATH, "../../div[contains(@class,'inventory-item-show-more')]/img").click()
    d.find_element(By.CSS_SELECTOR, ".inventory-modal-delete-btn img").click()

  add_item = d.find_element(
      By.CSS_SELECTOR, ".inventory-list-add-item-input input")

  add_item.send_keys("pepPEroni")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("cauliflower")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("caulIflowEr")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("BaCon")
  add_item.send_keys(Keys.TAB)

  assert len(d.find_elements(By.XPATH, inv_item_xpath % "pepperoni")) == 1
  assert len(d.find_elements(By.XPATH, inv_item_xpath % "cauliflower")) == 1
  assert len(d.find_elements(By.XPATH, inv_item_xpath % "bacon")) == 1

  item_ele = d.find_element(By.XPATH, inv_item_xpath % "bacon")
  clear(d, item_ele)
  item_ele.send_keys("pasta")
  item_ele.send_keys(Keys.TAB)
  assert len(w.until(lambda x: x.find_elements(
      By.XPATH, inv_item_xpath % "pasta"))) == 1
  assert len(d.find_elements(By.XPATH, inv_item_xpath % "bacon")) == 0

  item_ele = d.find_element(By.XPATH, inv_item_xpath % "cauliflower")
  clear(d, item_ele)
  item_ele.send_keys("strawberries")
  item_ele.send_keys(Keys.TAB)
  assert len(w.until(lambda x: x.find_elements(
      By.XPATH, inv_item_xpath % "strawberries"))) == 1
  assert len(d.find_elements(By.XPATH, inv_item_xpath % "cauliflower")) == 0

  d.find_element(
      By.XPATH, (inv_item_xpath % "pepperoni") +
      "/../../div[contains(@class,'inventory-item-show-more')]/img").click()
  d.find_element(By.CSS_SELECTOR, ".inventory-modal-delete-btn img").click()

  d.find_element(
      By.XPATH, (inv_item_xpath % "pasta") +
      "/../../div[contains(@class,'inventory-item-show-more')]/img").click()
  d.find_element(By.CSS_SELECTOR, ".inventory-modal-delete-btn img").click()

  d.find_element(
      By.XPATH, (inv_item_xpath % "strawberries") +
      "/../../div[contains(@class,'inventory-item-show-more')]/img").click()
  d.find_element(By.CSS_SELECTOR, ".inventory-modal-delete-btn img").click()

  assert len(d.find_elements(By.XPATH, inv_item_xpath % "pepperoni")) == 0
  assert len(d.find_elements(By.XPATH, inv_item_xpath % "pasta")) == 0
  assert len(d.find_elements(By.XPATH, inv_item_xpath % "strawberries")) == 0

  d.refresh()
  add_item = w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".inventory-list-add-item-input input"))
  assert len(d.find_elements(By.XPATH, inv_item_xpath % "pepperoni")) == 0
  assert len(d.find_elements(By.XPATH, inv_item_xpath % "pasta")) == 0
  assert len(d.find_elements(By.XPATH, inv_item_xpath % "strawberries")) == 0


def test_category_add():
  d, w = login()

  # Go to inventory list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-inventory img").click()

  add_item = d.find_element(
      By.CSS_SELECTOR, ".inventory-list-add-item-input input")

  add_item.send_keys("white wine")
  add_item.send_keys(Keys.ENTER)
  w.until(lambda x: x.find_element(By.XPATH, inv_item_xpath % "white wine"))

  d.find_element(By.ID, "inventory-list-cooking-btn").click()
  add_item.send_keys("zucchini")
  add_item.send_keys(Keys.ENTER)
  w.until(lambda x: x.find_element(By.XPATH, inv_item_xpath % "zucchini"))

  d.find_element(By.ID, "inventory-list-spices-btn").click()
  add_item.send_keys("cinnamon")
  add_item.send_keys(Keys.ENTER)
  w.until(lambda x: x.find_element(By.XPATH, inv_item_xpath % "cinnamon"))

  d.find_element(By.ID, "inventory-list-sauces-btn").click()
  add_item.send_keys("sriracha")
  add_item.send_keys(Keys.ENTER)
  w.until(lambda x: x.find_element(By.XPATH, inv_item_xpath % "sriracha"))

  d.find_element(By.ID, "inventory-list-non-cooking-btn").click()
  add_item.send_keys("tuna")
  add_item.send_keys(Keys.ENTER)
  w.until(lambda x: x.find_element(By.XPATH, inv_item_xpath % "ice cream"))

  d.find_element(By.ID, "inventory-list-non-cooking-btn").click()
  add_item.send_keys("ice cream")
  add_item.send_keys(Keys.ENTER)
  w.until(lambda x: x.find_element(By.XPATH, inv_item_xpath % "ice cream"))


def test_category_update():
  d, w = login()

  # Go to inventory list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-inventory img").click()
  d.find_element(By.ID, "inventory-list-cooking-btn").click()

  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".inventory-item-show-more img")).click()

  item = d.find_element(
      By.CSS_SELECTOR, ".inventory-modal-item input").get_attribute("value")

  categories = [
      ("spices", "Spices"),
      ("sauces", "Sauces and Condiments"),
      ("non-cooking", "Non-Cooking"),
      ("non-perishables", "Non-Perishables"),
      ("cooking", "Cooking"),
  ]
  for short_name, long_name in categories:
    # Move to spices
    category_ele = Select(d.find_element(
        By.CSS_SELECTOR, ".inventory-modal-category select"))
    category_ele.select_by_value(short_name)
    assert category_ele.first_selected_option.text == long_name
    d.find_element(By.CSS_SELECTOR, ".modal-panel-close-btn img").click()
    assert len(d.find_elements(
        By.XPATH, f"//div[contains(@class, 'inventory-item-input')]/input[@value='{item}']")) == 0

    # Find in spices
    d.find_element(
        By.ID, f"inventory-list-{short_name}-btn").click()
    item_input = w.until(lambda x: x.find_element(
        By.XPATH, f"//div[contains(@class, 'inventory-item-input')]/input[@value='{item}']"
    ))
    item_input.find_element(
        By.XPATH, "../../div[contains(@class, 'inventory-item-show-more')]/img").click()
    assert item == d.find_element(
        By.CSS_SELECTOR, ".inventory-modal-item input").get_attribute("value")
