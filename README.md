# REX2ANSI #

Convert [REXPaint](http://www.gridsagegames.com/rexpaint/) files to 24-bit color ANSI art.

Twenty. Four. Bit.  Color ANSI.

What's that, you say, 24-bit color ANSI art?  You out of touch human. ANSI art died when [Codepage 437](https://en.wikipedia.org/wiki/Code_page_437) faded everywhere except the Windows command prompt which doesn't even do ANSI color anymore.

Not true anymore!  This beautiful anachronism exists on most terminals!  Windows command prompt [recently gained](https://blogs.msdn.microsoft.com/commandline/2016/09/22/24-bit-color-in-the-windows-console/) support for 24-bit color here.    

Windows, Mac and Linux all have 24-bit ANSI art support in their various terminals.  REXPaint is one of the few programs that generates high color art.  All that remained was a way to get it out.

## Features ##

 - Standalone executable with CLI options, graphics pipeline tool
 - Export codepage 437 (classic ANSI) or UTF-8 (non-Windows users want this)
 - Supports layer flattening

## Screenshot ##

![Rex2ansi in action](/screenshots/1.jpg?raw=true)

## Usage ##

 1. `go build github.com/mlabbe/rex2ansi`
 
 2. `./rex2ansi --help`
