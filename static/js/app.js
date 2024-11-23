function Prompt() {
    const toast = function (toastObj) {
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = toastObj;
        const Toast = Swal.mixin({
            toast: true,
            icon: icon,
            position: position,
            title: msg,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.onmouseenter = Swal.stopTimer;
                toast.onmouseleave = Swal.resumeTimer;
            }
        });
        Toast.fire({});
    }
    const success = function (successObj) {
        const {
            msg = "",
            text = '',
            footer = "",
        } = successObj;


        Swal.fire({
            icon: "success",
            title: msg,
            text: text,
            footer: footer,
        });
    }
    const error = function (errorObj) {
        const {
            msg = "",
            text = '',
            footer = "",
        } = errorObj;


        Swal.fire({
            icon: "error",
            title: msg,
            text: text,
            footer: footer,
        });
    }
    const custom = async function (c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c;
        const {value: result} = await Swal.fire({
            icon:icon,
            title: title,
            html: msg,
            focusConfirm: false,
            backdrop: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen()
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen()
                }

            }
        });
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result);


                    }
                } else {
                    c.callback(false);
                }

            } else {
                c.callback(false);
            }
        }

    }
    return {
        toast: toast,
        success: success,
        error: error,
        custom
    }
}