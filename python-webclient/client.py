import random
import webclient_helper as wh
from pprint import pprint


def show_response(response):
    print(f'({response.status_code}) {response.request.method} {response.url}')
    try:
        pprint(response.json())
    except:
        print(response.content)
    print()


class ReceiptClient(wh.WebClient):
    uuids = []

    def post_and_save(self, json=None, debug=False, verbose=False):
        response = self.POST(
            '/receipts/process',
            json=json,
            debug=debug
        )
        if response.status_code == 200:
            self.uuids.append(response.json()['id'])
        if verbose:
            show_response(response)
        return response

    def get_points(self, _id=None, debug=False, verbose=False):
        if not _id and self.uuids:
            _id = random.choice(self.uuids)
            if verbose:
                print(f'Selected {_id=}')

        response = self.GET(
            f'/receipts/{_id}/points',
            debug=debug
        )
        if verbose:
            show_response(response)
        return response

    def get_breakdown(self, _id=None, debug=False, verbose=False):
        if not _id and self.uuids:
            _id = random.choice(self.uuids)
            if verbose:
                print(f'Selected {_id=}')

        response = self.GET(
            f'/receipts/{_id}/breakdown',
            debug=debug
        )
        if verbose:
            show_response(response)
        return response

    def get_bad1(self, debug=False, verbose=False):
        response = self.GET(
            '/receipts/process',
            debug=debug
        )
        if verbose:
            show_response(response)
        return response

    def get_bad2(self, debug=False, verbose=False):
        response = self.GET(
            '/receipts/bad',
            debug=debug
        )
        if verbose:
            show_response(response)
        return response

    def post_receipt1(self, debug=False, verbose=False):
        response = self.post_and_save(
            json={
                'retailer': 'Walgreens',
                'purchaseDate': '2022-01-02',
                'purchaseTime': '08:13',
                'items': [
                    {'shortDescription': 'Pepsi - 12-oz', 'price': '1.25'},
                    {'shortDescription': 'Dasani', 'price': '1.40'},
                ],
                'total': '2.65'
            },
            debug=debug,
            verbose=verbose
        )
        return response

    def post_receipt2(self, debug=False, verbose=False):
        response = self.post_and_save(
            json={
                'retailer': 'Target',
                'purchaseDate': '2022-01-01',
                'purchaseTime': '13:01',
                'items': [
                    {'shortDescription': 'Mountain Dew 12PK', 'price': '6.49'},
                    {'shortDescription': 'Emils Cheese Pizza', 'price': '12.25'},
                    {'shortDescription': 'Knorr Creamy Chicken', 'price': '1.26'},
                    {'shortDescription': 'Doritos Nacho Cheese', 'price': '3.35'},
                    {'shortDescription': '   Klarbrunn 12-PK 12 FL OZ  ', 'price': '12.00'}
                ],
                'total': '35.35'
            },
            debug=debug,
            verbose=verbose
        )
        return response

    def post_receipt3(self, debug=False, verbose=False):
        response = self.post_and_save(
            json={
                'retailer': 'M&M Corner Market',
                'purchaseDate': '2022-03-20',
                'purchaseTime': '14:33',
                'items': [
                    {'shortDescription': 'Gatorade', 'price': '2.25'},
                    {'shortDescription': 'Gatorade', 'price': '2.25'},
                    {'shortDescription': 'Gatorade', 'price': '2.25'},
                    {'shortDescription': 'Gatorade', 'price': '2.25'}
                ],
                'total': '9.00'
            },
            debug=debug,
            verbose=verbose
        )
        return response

    def post_bad1(self, debug=False, verbose=False):
        response = self.post_and_save(
            json={
                'retailer': 'M&M Corner Market',
                'purchaseDate': '2022-13-20',
                'purchaseTime': '24:33',
                'items': [
                    {'shortDescription': 'Stuff', 'price': '2.25'},
                    {'shortDescription': 'Stuff', 'price': '2.25'},
                ],
                'total': '4.50'
            },
            debug=debug,
            verbose=verbose
        )
        return response

    def post_bad2(self, debug=False, verbose=False):
        response = self.post_and_save(
            json={
                'retailer': 'M&M Corner Market',
                'purchaseDate': '2022-03-20',
                'purchaseTime': '14:33',
                'items': [
                    {'shortDescription': 'Stuff', 'price': '2.25'},
                    {'shortDescription': 'Stuff', 'price': '2.25'},
                ],
                'total': '4.60'
            },
            debug=debug,
            verbose=verbose
        )
        return response

    def get_all_points_and_breakdowns(self, debug=False, verbose=False):
        responses = []
        for _id in self.uuids:
            responses.append(self.get_points(_id=_id, debug=debug, verbose=verbose))
            responses.append(self.get_breakdown(_id=_id, debug=debug, verbose=verbose))
        return responses


client = ReceiptClient(
    base_url='http://localhost:8080'
)


if __name__ == '__main__':
    client.get_points(verbose=True)
    client.get_bad1(verbose=True)
    client.get_bad2(verbose=True)
    client.post_receipt1(verbose=True)
    client.post_receipt2(verbose=True)
    client.post_receipt3(verbose=True)
    client.post_bad1(verbose=True)
    client.post_bad2(verbose=True)
    client.get_points(verbose=True)
    client.get_breakdown(verbose=True)
    client.get_points('abc-123', verbose=True)
    client.get_breakdown('abc-123', verbose=True)
    client.get_all_points_and_breakdowns(verbose=True)
