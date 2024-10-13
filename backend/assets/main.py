import json

def update_card_ids(filename):
    # Load the existing cards from JSON file
    with open(filename, 'r') as file:
        cards = json.load(file)

    # Update the ID for each card
    for i, card in enumerate(cards):
        card['id'] = i + 1  # Start IDs from 1

    # Save the updated cards back to the JSON file
    with open(filename, 'w') as file:
        json.dump(cards, file, indent=4)

if __name__ == "__main__":
    update_card_ids('cards.json')

