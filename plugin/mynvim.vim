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


function! mynvim#begin(plugin_dir) abort
    let g:mynvim_plugin_dir = expand(a:plugin_dir)
endfunction

" -----------------------------------------------------------------------------
" plugin manifest

call remote#host#Register(s:plugin_name, '', function('s:JobStart'))

call remote#host#RegisterPlugin('mynvim', '0', [
\ {'type': 'command', 'name': 'Expand', 'sync': 0, 'opts': {'range': ''}},
\ {'type': 'command', 'name': 'Plugin', 'sync': 0, 'opts': {'nargs': '+'}},
\ {'type': 'command', 'name': 'PluginDir', 'sync': 0, 'opts': {'nargs': '+'}},
\ {'type': 'command', 'name': 'PluginInstall', 'sync': 0, 'opts': {'nargs': '0'}},
\ {'type': 'command', 'name': 'Push', 'sync': 0, 'opts': {'nargs': '+'}},
\ {'type': 'function', 'name': 'PushBuf', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'ToggleContainer', 'sync': 1, 'opts': {}},
\ ])
