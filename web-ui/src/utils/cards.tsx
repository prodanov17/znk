// src/helpers/cards.tsx
const cardImages = import.meta.glob('../assets/*.svg');

// Function to generate the card name based on suit and rank
const getCardName = (suit: string, rank: string): string => {
    const suitAbbreviation = suit.charAt(0).toLowerCase(); // Get the first letter of the suit
    const rankAbbreviation = rank === '10' ? '10' : rank.charAt(0).toLowerCase(); // Handle '10' separately
    return `${rankAbbreviation}${suitAbbreviation}`; // e.g., '10c', 'ac', etc.
};

// Function to get the card image based on suit and rank
export const getCardImage = async (suit: string, rank: string): Promise<string> => {
    const cardName = getCardName(suit, rank);
    const cardImage = cardImages[`../assets/${cardName}.svg`];

    if (cardImage) {
        const imageModule: any = await cardImage(); // Dynamically import the SVG
        return imageModule.default; // Return the default export from the SVG
    }

    throw new Error(`Card image not found for ${rank} of ${suit}`);
};

