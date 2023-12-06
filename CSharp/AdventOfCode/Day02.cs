namespace AdventOfCode;

readonly record struct Game(int ID, Dictionary<string, int>[] Data);

public sealed class Day02 : BaseDay
{
    private readonly List<Game> _games;

    public Day02()
    {
        _games = new List<Game>();
        foreach (var line in File.ReadLines(InputFilePath))
        {
            var parts = line.Split(": ");
            int id = int.Parse(parts[0].Split(" ")[1]);

            var gameDataText = parts[1].Split("; ");
            var gameData = gameDataText.Select(ParseGameData).ToArray();
            _games.Add(new Game{Data = gameData, ID = id});
        }
    }

    public override ValueTask<string> Solve_1()
    {
        var contents = new Dictionary<string, int>()
        {
            ["red"] = 12,
            ["green"] = 13,
            ["blue"] = 14
        };

        return new(
            _games
                .Where(g => g.Data.All(d => d.All(pair => contents[pair.Key] >= pair.Value)))
                .Sum(g => g.ID)
                .ToString());
    }

    public override ValueTask<string> Solve_2()
    {
        var setsRequired = _games.Select(g =>
        {
            var maxes = new Dictionary<string, int>()
            {
                ["red"] = 0,
                ["green"] = 0,
                ["blue"] = 0
            };
            foreach (var dict in g.Data)
            {
                foreach (var (k, v) in dict)
                    if (maxes[k] < v)
                        maxes[k] = v;
            }

            return maxes;
        });

        return new(setsRequired.Select(PowerSetOf).Sum().ToString());
    }

    private int PowerSetOf(Dictionary<string, int> dict)
        => dict.Select(pair => pair.Value).Aggregate((a, b) => a * b);

    private Dictionary<string, int> ParseGameData(string s)
    {
        var parts = s.Split(", ");
        return parts.Select(ParseColorPair).ToDictionary();
    }

    private (string, int) ParseColorPair(string s)
    {
        var parts = s.Split(" ");
        return (parts[1], int.Parse(parts[0]));
    }
}
