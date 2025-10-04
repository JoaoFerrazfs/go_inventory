// src/pages/Racks.jsx
import React, { useEffect, useState } from "react";

function Racks() {
    const [racks, setRacks] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetch("http://localhost:3000/api/v1/racks/")

            .then((res) => {
                if (!res.ok) throw new Error(`Erro HTTP: ${res.status}`);
                return res.json();
            })
            .then((data) => {
                setRacks(data);
                setLoading(false);
            })
            .catch((err) => {
                setError(err.message);
                setLoading(false);
            });
    }, []);

    const renderProgressBar = (percentage) => {
        let color = "bg-green-500";
        if (percentage > 80) color = "bg-red-500";
        else if (percentage > 50) color = "bg-yellow-500";
        return (
            <div className="w-full bg-gray-200 rounded-full h-2.5">
                <div
                    className={`h-2.5 rounded-full ${color}`}
                    style={{ width: `${percentage}%` }}
                />
            </div>
        );
    };

    return (
        <div className="container mx-auto p-4 flex-1">
            <header className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold text-blue-900">Racks Cadastrados ruim</h1>
                <a
                    href="/racks/novo"
                    className="bg-green-600 hover:bg-green-700 text-white font-semibold py-2 px-4 rounded-full transition-colors shadow-lg flex items-center"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-5 w-5 mr-1"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                    >
                        <path
                            fillRule="evenodd"
                            d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z"
                            clipRule="evenodd"
                        />
                    </svg>
                    Novo
                </a>
            </header>

            {loading && (
                <p className="col-span-full text-center text-gray-500">
                    Carregando racks...
                </p>
            )}

            {error && (
                <div className="col-span-full text-center p-10 bg-white rounded-xl shadow-lg mt-8 text-red-600">
                    <p className="text-lg">Erro ao carregar os dados. Tente novamente.</p>
                    <p className="text-sm">Detalhe: {error}</p>
                </div>
            )}

            {!loading && !error && racks.length === 0 && (
                <div className="col-span-full text-center p-10 bg-white rounded-xl shadow-lg mt-8">
                    <p className="text-gray-500 text-lg">
                        Nenhum rack cadastrado ainda. Comece adicionando um novo! 🚀
                    </p>
                </div>
            )}

            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                {!loading &&
                    !error &&
                    racks.map((rack) => {
                        const slotsOcupados = rack.Pallets ? rack.Pallets.length : 0;
                        const percentage = parseFloat(rack.PercetageUsed || 0).toFixed(0);

                        return (
                            <a
                                key={rack.ID}
                                href={`/racks/${rack.ID}`}
                                className="block rack-link"
                            >
                                <div className="bg-white p-4 rounded-xl shadow-lg hover:shadow-xl transition-shadow border-t-4 border-blue-600">
                                    <div className="flex items-center justify-center mb-4">
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            className="h-16 w-16 text-blue-500"
                                            viewBox="0 0 20 20"
                                            fill="currentColor"
                                        >
                                            <path
                                                fillRule="evenodd"
                                                d="M3 3a1 1 0 00-1 1v12a1 1 0 001 1h14a1 1 0 001-1V4a1 1 0 00-1-1H3zm1 2h12v2H4V5zm0 4h12v2H4V9zm0 4h12v2H4v-2z"
                                                clipRule="evenodd"
                                            />
                                        </svg>
                                    </div>

                                    <h2
                                        className="text-xl font-semibold text-gray-800 mb-2 truncate rack-name"
                                        title={rack.Name}
                                    >
                                        {rack.Name}
                                    </h2>

                                    <p className="text-sm text-gray-500 mb-1">
                                        <strong>Localização:</strong>{" "}
                                        <span className="font-medium text-gray-700">
                                            {rack.Location}
                                        </span>
                                    </p>

                                    <p className="text-sm text-gray-500 mb-3">
                                        <strong>Capacidade:</strong>{" "}
                                        <span className="font-medium text-gray-700">
                                            {slotsOcupados} / {rack.TotalCapacity} Pallets
                                        </span>
                                    </p>

                                    {renderProgressBar(percentage)}

                                    <p className="text-xs text-right mt-1 text-gray-600">
                                        {percentage}% Ocupado
                                    </p>
                                </div>
                            </a>
                        );
                    })}
            </div>
        </div>
    );
}

export default Racks;
