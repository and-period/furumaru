import { setActivePinia, createPinia } from 'pinia'

import { AddSnackbarPayload, useCommonStore } from '~~/src/store/common'

describe('Common Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('default: snackbars state is empty list', () => {
    const { snackbars } = useCommonStore()
    expect(snackbars).toEqual([])
  })

  it('when addSnackBar is called, state is updated', () => {
    const { snackbars, addSnackbar } = useCommonStore()
    const snackbarPayload: AddSnackbarPayload = {
      message: 'ようこそ、ふるまるへ',
      color: 'info',
    }
    addSnackbar(snackbarPayload)

    expect(snackbars.length).toBe(1)
    expect(snackbars[0].message).toEqual(snackbarPayload.message)
    expect(snackbars[0].color).toEqual(snackbarPayload.color)
  })

  it('info addSnack default timeout value is 5000', () => {
    const { snackbars, addSnackbar } = useCommonStore()
    const snackbarPayload: AddSnackbarPayload = {
      message: 'ようこそ、ふるまるへ',
      color: 'info',
    }
    addSnackbar(snackbarPayload)

    expect(snackbars.length).toBe(1)
    expect(snackbars[0].timeout).toBe(5000)
  })

  it('error addSnack default timeout value is -1', () => {
    const { snackbars, addSnackbar } = useCommonStore()
    const snackbarPayload: AddSnackbarPayload = {
      message: 'エラーが発生しました',
      color: 'error',
    }
    addSnackbar(snackbarPayload)

    expect(snackbars.length).toBe(1)
    expect(snackbars[0].timeout).toBe(-1)
  })

  it('when hideSnackbar called,  state is updated', () => {
    const { snackbars, addSnackbar, hideSnackbar } = useCommonStore()
    addSnackbar({
      message: 'ようこそ、ふるまるへ',
      color: 'info',
    })

    expect(snackbars.length).toBe(1)

    hideSnackbar(0)
    expect(snackbars.length).toBe(0)
  })
})
