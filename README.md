A build pre-processor tool for DayZ mods.  This is a work in progress.

# Usage

```
Usage of mod-build:
  -config string
        config file (optional)
  -image-to-paa string
        Path to the ImageToPAA executable (default "C:\\Program Files (x86)\\Steam\\steamapps\\common\\DayZ Tools\\Bin\\ImageToPAA\\ImageToPAA.exe")
  -output string
        Path to the output directory (default "./build/")
  -source string
        Path to the source directory (default "./source/")
  -yes
        Automatically confirm all prompts (use with caution)
```

## Example Output

```
===================================================
ImageToPAA Path: ImageToPAA.exe
    Source Path: ./source/
    Output Path: ./build/
   Auto-confirm: false
===================================================
⚠️ The build directory "build/WILDLANDZ_Core" will be removed and recreated. Continue? [y/N] n
⚠️ The build directory "build/WILDLANDZ_GreenCounty" will be removed and recreated. Continue? [y/N] y
⏭️ Skipping: "WILDLANDZ_Core"
🛠️ Building: WILDLANDZ_GreenCounty
   📂 Creating build output directory "build/WILDLANDZ_GreenCounty"
   📄 Copying    : "$PBOPREFIX@.txt"
   📄 Copying    : "config.cpp"
   📄 Copying    : "weapons/attachments/magazine/config.cpp"
   🔁 Converting : "weapons/attachments/magazine/data/pmag_gc_co.png"
   🔁 Converting : "weapons/firearms/ak101/data/ak101_gc_co.png"
   🔁 Converting : "weapons/firearms/ak74u/data/aks74u_gc_co.png"
   📄 Copying    : "weapons/firearms/akm/config.cpp"
   🔁 Converting : "weapons/firearms/akm/data/akm_gc_co.png"
   🔁 Converting : "weapons/firearms/izh18/data/izh18_gc_co.png"
   🔁 Converting : "weapons/firearms/m16a2/data/m16a2_gc_co.png"
   📄 Copying    : "weapons/firearms/m4/config.cpp"
   🔁 Converting : "weapons/firearms/m4/data/m4_body_gc_co.png"
🎉 Done!
```
