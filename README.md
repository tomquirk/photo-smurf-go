# Photo Smurf Go

### Batch Photo Library Cleanup in Go

## Why?

> All my photos are unorganised, scattered across my computer/harddrive - help!

Sound familiar?

Simply tell these smurfs when you went on holidays i.e New Zealand between 27/01/2016 and 15/02/2016. Photo smurf will search your computer for any photo taken between specified dates and put them in folders representing albums!

This project is a Go implementation of [photo-smurf](https://github.com/tomquirk/photo-smurf), originally written in Python. I haven't done benchmarks, but `Go > Python` :joy:

## Usage

Create an _albums_ config file as described below

### CLI
```
go install github.com/tomquirk/photosmurf
photosmurf [srcRootPath] [destRootPath] [album_config]
```

## The Album Config File
- Make sure your album name contains NO SPACES (because that's silly)!

#### Example
> albums.txt

```
[
  {
    "name": "p_holiday_tasmania-2016",
    "startTime": "18 Dec 16 00:00 UTC",
    "endTime": "26 Dec 16 00:00 UTC"
  },
  ...
]
```

# Todo
-   tests
-   docs/general cleanup
