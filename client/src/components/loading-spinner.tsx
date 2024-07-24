export default function LoadingSpinner() {
  return (
    <div className="fixed inset-0 bg-gray-700 bg-opacity-50 flex justify-center items-center">
      <div
        className="inline-block h-16 w-16 animate-spin rounded-full border-4 border-solid border-current border-t-transparent align-[-0.125em] text-surface motion-reduce:animate-[spin_1.5s_cubic_bezier(0.4,0,0.6,1)_infinite] text-slate-500"
        role="status"
      >
        {/* <span className="visually-hidden">Loading...</span> */}
      </div>
    </div>
  );
}
