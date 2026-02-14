A build pre-processor tool for DayZ mods.  This is a work in progress.

# Features

- Copies known file types to output directory
- Converts .png or .jpg files to .paa
- Only copies/convert files that have changed

# Usage

```
usage: mod-build [options] <source-directory>
  -clean
        Clean output directory before building (deletes files which are not present in the source)
  -config string
        config file (optional)
  -image-to-paa string
        Path to the ImageToPAA executable (default "C:\\Program Files (x86)\\Steam\\steamapps\\common\\DayZ Tools\\Bin\\ImageToPAA\\ImageToPAA.exe")
  -output string
        Path to the output directory root (where built addons will be placed) (default "P:\\")
  -yes
        Automatically confirm all prompts (use with caution)
```

## Example Output

```
===================================================
ImageToPAA Path: ImageToPAA.exe
    Source Path: source/WILDLANDZ_Anniversary
    Output Root: P:/
   Auto-confirm: false
          Clean: false
---------------------------------------------------
      Addon Name: WILDLANDZ_Anniversary
Output Directory: P:/WILDLANDZ_Anniversary
===================================================
⚠️ The contents of "P:/WILDLANDZ_Anniversary" will be removed or replaced. Continue? [y/N] y
⏭️ Unchanged  : "$PBOPREFIX@.txt"
⏭️ Unchanged  : "characters/backpacks/config.cpp"
📄 Copying    : "config.cpp"
⏭️ Unchanged  : "gear/consumables/config.cpp"
⏭️ Unchanged  : "gear/food/config.cpp"
⏭️ Unchanged  : "gear/food/data/cupcake.p3d"
⏭️ Unchanged  : "gear/food/data/cupcake.rvmat"
⏭️ Unchanged  : "gear/food/data/cupcake_rotten.rvmat"
⏭️ Unchanged  : "characters/backpacks/data/armypouch_anniversary_co.png"
⏭️ Unchanged  : "gear/consumables/data/anniversary_paper1_1.png"
⏭️ Unchanged  : "gear/consumables/data/anniversary_paper1_2.png"
⏭️ Unchanged  : "gear/consumables/data/anniversary_paper1_3.png"
⏭️ Unchanged  : "gear/consumables/data/anniversary_ribbon_2.png"
⏭️ Unchanged  : "gear/consumables/data/anniversary_ribbon_3.png"
⏭️ Unchanged  : "gear/consumables/data/anniversary_ribbon_co1.png"
🔁 Converting : "gear/food/data/cupcake.png"
⏭️ Unchanged  : "gear/food/data/cupcake_nohq.png"
⏭️ Unchanged  : "gear/food/data/cupcake_smdi.png"
🎉 Done!
```
