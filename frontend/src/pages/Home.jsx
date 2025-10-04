
export default function Home() {
    return (
        <div className="flex flex-col min-h-screen bg-gray-100 font-sans">


            <main className="flex flex-col items-center text-center flex-1 px-4 pt-24 md:pt-32">
                <h1 className="text-3xl md:text-4xl font-bold text-blue-900 mb-4">
                    Bem-vindo ao Go Inventory
                </h1>
                <p className="text-gray-600 max-w-md mb-6">
                    Um site de inventário logístico, pensado para simplificar o controle de
                    produtos e estoques de forma eficiente e acessível.
                </p>
                <img
                    src="/img/banner1.png"
                    alt="Banner"
                    className="w-4/5 max-w-lg rounded-lg mb-6 object-contain"
                />
                <a
                    href="/login"
                    className="bg-red-600 hover:bg-red-700 text-white font-semibold py-3 px-6 rounded transition-colors"
                >
                    Fazer Login
                </a>
            </main>

        </div>
    );
}
