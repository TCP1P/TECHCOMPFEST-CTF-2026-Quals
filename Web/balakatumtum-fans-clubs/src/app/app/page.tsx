import Image from "next/image";
import tongue from "../public/tongue.png";

export default function Home() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-yellow-50 font-sans dark:bg-zinc-900">
      <main className="flex min-h-screen w-full max-w-4xl flex-col items-center justify-between py-24 px-8 sm:items-center">
        
        {/* TOP SECTION: THE PIZZA */}
        <div className="relative flex flex-col items-center group">
          <div className="relative drop-shadow-2xl transition-transform hover:rotate-180 duration-700 ease-in-out">
            {/* Using a working Unsplash URL for the Pizza */}
            <img
              src="https://images.unsplash.com/photo-1513104890138-7c749659a591?auto=format&fit=crop&w=400&q=80"
              alt="Delicious Pizza"
              width={300}
              height={300}
              className="rounded-full border-4 border-white dark:border-zinc-800 object-cover h-[300px] w-[300px]"
            />
          </div>
          <h2 className="mt-8 text-6xl font-black tracking-tighter text-red-600 dark:text-red-500 uppercase drop-shadow-sm">
            The Pizza Guy
          </h2>
        </div>

        {/* MIDDLE SECTION: TEXT & TONGUE */}
        <div className="flex flex-col items-center gap-6 text-center">
          <p className="max-w-md text-2xl font-medium leading-8 text-zinc-700 dark:text-zinc-300">
            Cheesy. Crunchy. <br/>
            <span className="font-black text-black dark:text-white">
              Lick your screen good.
            </span>
          </p>

          {/* THE TONGUE IMAGE */}
          <div className="animate-bounce mt-4">
            {/* Using a working Wikimedia URL for the Tongue (Einstein style) */}
            <img
              src={tongue.src}
              alt="Tongue sticking out"
              width={120}
              height={150}
              className="object-cover rounded-xl rotate-12 shadow-xl border-4 border-white dark:border-zinc-800"
            />
          </div>
        </div>

        {/* BOTTOM SECTION: BUTTONS */}
        <div className="flex flex-col gap-4 text-base font-bold sm:flex-row mt-8">
          <button
            className="flex h-14 w-full items-center justify-center gap-2 rounded-full bg-red-600 px-10 text-white shadow-lg shadow-red-600/30 transition-all hover:bg-red-700 hover:scale-105 active:scale-95 md:w-auto"
          >
            🍕 GIMMIE PIZZA
          </button>
          <button
            className="flex h-14 w-full items-center justify-center rounded-full border-2 border-zinc-900 bg-white px-10 text-zinc-900 transition-all hover:bg-zinc-100 md:w-auto dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100"
          >
            See Menu
          </button>
        </div>
      </main>
    </div>
  );
}