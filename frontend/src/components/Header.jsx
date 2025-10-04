import React from "react";

export default function Header() {
    return (
        <header className="flex justify-between items-center bg-blue-900 text-white px-4 py-3">
            <div className="text-2xl font-bold">Go Inventory</div>
            <script src="https://cdn.tailwindcss.com"></script>
            <img
                src="/static/img/logo2.png"
                alt="Logo"
                className="h-12 md:h-14 object-contain"
            />
        </header>
    );
}
