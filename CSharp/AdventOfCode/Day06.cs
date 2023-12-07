namespace AdventOfCode;

public sealed class Day06 : BaseDay
{
    private readonly (long time, long recordDistance)[] _races;
    private readonly (long time, long recordDistance) _racesPart2;

    public Day06()
    {
        var lines = File.ReadAllLines(InputFilePath);
        
        var times = lines[0][("Time:".Length + 1)..].Split(' ', StringSplitOptions.RemoveEmptyEntries).Select(long.Parse).ToArray();
        var distances = lines[1][("Distance:".Length + 1)..].Split(' ', StringSplitOptions.RemoveEmptyEntries).Select(long.Parse).ToArray();

        _races = times.Zip(distances).ToArray();
        
        var time2 = long.Parse(lines[0][("Time:".Length + 1)..].Replace(" ", ""));
        var dist2 = long.Parse(lines[1][("Distance:".Length + 1)..].Replace(" ", ""));
        _racesPart2 = (time2, dist2);
    }

    public override ValueTask<string> Solve_1()
        => new(_races.Select(GetWaysToWin).Aggregate(1L, (a, i) => i * a).ToString());

    public override ValueTask<string> Solve_2() => new(GetWaysToWin(_racesPart2).ToString());

    private static long GetWaysToWin((long time, long recordDistance) race)
        => Enumerable.Range(0, (int)race.time).Aggregate(0, (acc, t) => t * (race.time - t) > race.recordDistance ? acc + 1 : acc);
}
