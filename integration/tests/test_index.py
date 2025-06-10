from tests import login_oauth, login_traditional


def test_login_traditional():
  login_traditional()


def test_login_oauth():
  login_oauth()
