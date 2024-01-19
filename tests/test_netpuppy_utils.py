from netpuppy.utils import sum


def test_sum() -> None:
    assert sum(1, 2) == 3
    assert sum(1, 1) == 2
