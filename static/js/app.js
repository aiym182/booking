function Prompt() {
  let toast = function (c) {
    const { msg = '', icon = 'success', position = 'top-end' } = c

    const Toast = Swal.mixin({
      toast: true,
      title: msg,
      position: position,
      icon: icon,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,

      didOpen: (toast) => {
        toast.addEventListener('mouseenter', Swal.stopTimer)
        toast.addEventListener('mouseleave', Swal.resumeTimer)
      },
    })

    Toast.fire()
  }
  let success = function (c) {
    const { icon = 'success', title = '', msg = '' } = c

    Swal.fire({
      icon: icon,
      title: title,
      text: msg,
    })
  }

  let error = (c) => {
    const { icon = 'error', title = '', msg = '' } = c

    Swal.fire({
      icon: icon,
      title: title,
      text: msg,
    })
  }

  async function custom(c) {
    const { icon = '', msg = '', title = '', showConfirmButton = true } = c

    const { value: result } = await Swal.fire({
      title: title,
      icon: icon,
      html: msg,
      backdrop: false,
      focusConfirm: false,
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
      },
    })

    if (result) {
      if (result.dismiss !== Swal.DismissReason.cancel) {
        if (result.value != '') {
          if (c.callback !== undefined) {
            c.callback(result)
          } else {
            c.callback(false)
          }
        } else {
          c.callback(false)
        }
      }
    }
  }

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  }
}
