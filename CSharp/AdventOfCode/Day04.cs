namespace AdventOfCode;

record struct Card(HashSet<int> WinningNumbers, HashSet<int> ChosenNumbers, int Copies = 1)
{
    public int GetNumMatches()
    {
        var hs = new HashSet<int>(ChosenNumbers);
        hs.IntersectWith(WinningNumbers);
        return hs.Count;
    }
}

// ReSharper disable once InconsistentNaming
public sealed class Day04 : BaseDay
{
    private readonly Card[] _cards;

    public Day04()
    {
        _cards = File.ReadLines(InputFilePath)
            .Select(line =>
            {
                var data = line[(line.IndexOf(": ", StringComparison.Ordinal) + 2)..];

                var numData = data.Split(" | ");
                var winningNumbers = numData[0].Split(" ", StringSplitOptions.RemoveEmptyEntries).Select(int.Parse)
                    .ToHashSet();
                var chosenNumbers = numData[1].Split(" ", StringSplitOptions.RemoveEmptyEntries).Select(int.Parse)
                    .ToHashSet();
                return new Card(winningNumbers, chosenNumbers);
            }).ToArray();
    }

    public override ValueTask<string> Solve_1()
    {
        var totalScore = _cards
            .Select(card => card.GetNumMatches())
            .Select(c => (int)Math.Pow(2, c - 1))
            .Sum();

        return new(totalScore.ToString());
    }

    public override ValueTask<string> Solve_2()
    {
        for (int i = 0; i < _cards.Length; i++)
        {
            int matches = _cards[i].GetNumMatches();
            for (int j = 1; j <= matches; j++)
            {
                _cards[i + j].Copies += _cards[i].Copies;
            }
        }
        
        return new(_cards.Select(i => i.Copies).Sum().ToString());
    }
}
