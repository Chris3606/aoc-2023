namespace AdventOfCode;

readonly record struct Range(long Start, long End)
{
    public static Range FromLength(long start, long length) => new(start, start + length - 1);
    
    public bool Contains(long value) => value >= Start && value <= End;

    public bool Overlaps(Range other) => Start <= other.End && other.Start <= End;
}

readonly record struct MapEntry(Range Source, long DestinationStart);

public sealed class Day05 : BaseDay
{
    private readonly long[] _seedNumbers;
    private readonly MapEntry[][] _maps;

    public Day05()
    {
        var groups = File.ReadAllText(InputFilePath).Split(Environment.NewLine + Environment.NewLine);
        _seedNumbers = groups[0][(("seeds:".Length) + 1)..].Split(" ", StringSplitOptions.RemoveEmptyEntries).Select(long.Parse).ToArray();

        _maps = groups[1..].Select(i =>
        {
            var lines = i.Split(Environment.NewLine);
            return lines[1..].Select(j =>
            {
                var map = j.Split(" ", StringSplitOptions.RemoveEmptyEntries);
                var source = Range.FromLength(long.Parse(map[1]), long.Parse(map[2]));
                var destinationStart = long.Parse(map[0]);

                return new MapEntry(source, destinationStart);
            }).ToArray();
        }).ToArray();
    }

    public override ValueTask<string> Solve_1()
    {
        return new(_seedNumbers.Select(i => ApplyMaps(i, _maps)).Min().ToString());
    }

    public override ValueTask<string> Solve_2()
    {
        var seeds = _seedNumbers.Chunk(2).Select(i => Range.FromLength(i[0], i[1])).ToList();
        
        foreach (var map in _maps)
            ApplyMap(seeds, map);
        
        return new(seeds.Select(i => i.Start).Min().ToString());
    }

    private static void ApplyMap(List<Range> seeds, MapEntry[] map)
    {
        var shiftDeltas = new List<long>();
        for (int i = 0; i < seeds.Count; i++)
            shiftDeltas.Add(0);
    
        foreach (var range in map)
        {
            int len = seeds.Count;
            
            for (int i = 0; i < len; i++)
            {
                var seed = seeds[i];
                if (range.Source.Overlaps(seed))
                {
                    var overlap = new Range(Math.Max(seed.Start, range.Source.Start), Math.Min(seed.End, range.Source.End));

                    if (overlap.End < seed.End)
                    {
                        seeds.Add(seed with { Start = overlap.End + 1 });
                        shiftDeltas.Add(0);
                    }

                    if (overlap.Start > seed.Start)
                    {
                        seeds.Add(seed with { End = overlap.Start - 1 });
                        shiftDeltas.Add(0);
                    }

                    shiftDeltas[i] = range.DestinationStart - range.Source.Start;

                    seeds[i] = overlap;
                }
            }
        }
        
        for (int i = 0; i < seeds.Count; i++)
            seeds[i] = new Range(Start: seeds[i].Start + shiftDeltas[i], End: seeds[i].End + shiftDeltas[i]);
    }

    private static long ApplyMaps(long seedValue, MapEntry[][] maps)
    {
        foreach (var map in maps)
            foreach (var mapEntry in map)
            {
                if (mapEntry.Source.Contains(seedValue))
                {
                    seedValue = mapEntry.DestinationStart + (seedValue - mapEntry.Source.Start);
                    break;
                }
            }

        return seedValue;
    }
}
