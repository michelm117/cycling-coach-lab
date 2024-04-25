const toasts = [];
let isToastQueueRunning = false;

const onMakeToast = async (e) => {
  const level = e.detail.level;
  const message = e.detail.message;

  // Add toast to the queue
  //toasts.push({ level, message });
  window.toast(message, { type: level });

  // If the toast queue is not running, start it
  if (!isToastQueueRunning) {
    isToastQueueRunning = true;
    runToastQueue();
  }
};

const runToastQueue = async () => {
  while (toasts.length > 0) {
    const { level, message } = toasts.shift(); // Get the first toast from the queue

    const Toast = Swal.mixin({
      toast: true,
      position: "top-right",
      iconColor: "white",
      customClass: {
        popup: "colored-toast",
      },
      timer: 5000,
      showConfirmButton: false,
      timerProgressBar: true,
    });

    await Toast.fire({
      icon: level,
      title: message,
    });

    // Delay between toasts
    await new Promise((resolve) => setTimeout(resolve, 500));
  }

  // No more toasts in the queue, reset the flag
  isToastQueueRunning = false;
};

document.body.addEventListener("showToast", onMakeToast);

// simulate to test
// window.onload = async function () {
//   onMakeToast({
//     detail: {
//       level: "info",
//       message: "Welcome to the page!",
//     },
//   });
//   onMakeToast({
//     detail: {
//       level: "info",
//       message: "Welcome to the page!",
//     },
//   });
// };
