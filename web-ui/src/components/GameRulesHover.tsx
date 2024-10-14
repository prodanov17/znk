import { useState } from "react";
import { InfoIcon } from "./Icons";

type GameRulesProps = {
    rules: Record<string, string>;
};

const GameRulesHover = ({ rules }: GameRulesProps) => {
    const [showRules, setShowRules] = useState(false);

    return (
        <div
            className="absolute left-4 top-4 inline-block"
            onMouseEnter={() => setShowRules(true)}
            onMouseLeave={() => setShowRules(false)}
        >
            <InfoIcon className="text-neutral-400 cursor-pointer" />
            {showRules && (
                <div className="absolute left-0 mt-2 w-64 p-2 bg-gray-800 text-white rounded shadow-lg z-10">
                    <h4 className="font-semibold mb-2">Game Rules</h4>
                    <ul className="text-sm">
                        {Object.entries(rules).map(([key, value], index) => (
                            <div key={index}>
                                <strong>{key}:</strong> {value}
                            </div>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
};

export default GameRulesHover;

