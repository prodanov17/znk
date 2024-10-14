// src/components/Card.tsx
import { useEffect, useState } from 'react';
import { getCardImage } from '../utils/cards';

interface CardProps {
    suit: string;
    rank: string;
    className?: string;
    style?: React.CSSProperties;
}

const Card = ({ suit, rank, className = "", style }: CardProps) => {
    const [cardImage, setCardImage] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchCardImage = async () => {
            try {
                const image = await getCardImage(suit.charAt(0), rank.toLowerCase());
                setCardImage(image);
            } catch (err) {
                setError((err as Error).message);
            }
        };

        fetchCardImage();
    }, [suit, rank]);

    if (error) {
        return <div className="error">{error}</div>;
    }

    return (
        <div className={className} style={style}>
            {cardImage ? <img src={cardImage} alt={`${rank} of ${suit}`} className="w-full rounded-md" /> : <p>Loading...</p>}
        </div>
    );
};

export default Card;

