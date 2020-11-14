if exists('g:loaded_mynvim')
  finish
endif
let g:loaded_mynvim = 1

" -----------------------------------------------------------------------------
" register remote plugin

let s:plugin_name   = 'mynvim'
let s:plugin_root   = fnamemodify(resolve(expand('<sfile>:p')), ':h:h')

let s:plugin_cmd = [s:plugin_root . '/bin/' . s:plugin_name]

function! s:JobStart(host) abort
    return jobstart(s:plugin_cmd, {'rpc': v:true, 'detach': v:false})
endfunction

" -----------------------------------------------------------------------------
" plugin manifest

call remote#host#Register(s:plugin_name, '', function('s:JobStart'))

call remote#host#RegisterPlugin('mynvim', '0', [
\ {'type': 'function', 'name': 'NewContainer', 'sync': 1, 'opts': {}},
\ ])
