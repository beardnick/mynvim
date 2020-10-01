if exists('g:loaded_win_container')
  finish
endif
let g:loaded_win_container = 1

" -----------------------------------------------------------------------------
" register remote plugin

let s:plugin_name   = 'win-container'
let s:plugin_root   = fnamemodify(resolve(expand('<sfile>:p')), ':h:h')

let s:plugin_cmd = [s:plugin_root . '/bin/' . s:plugin_name]

function! s:JobStart(host) abort
    return jobstart(s:plugin_cmd, {'rpc': v:true, 'detach': v:false})
endfunction

" -----------------------------------------------------------------------------
" plugin manifest

call remote#host#Register(s:plugin_name, '', function('s:JobStart'))

call remote#host#RegisterPlugin('win-container', '0', [
\ {'type': 'function', 'name': 'NewContainer', 'sync': 1, 'opts': {}},
\ ])
