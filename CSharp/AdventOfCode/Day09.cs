namespace AdventOfCode;

public sealed class Day09 : BaseDay
{
    private readonly int[][] _sequences;

    public Day09()
    {
        _sequences = File.ReadLines(InputFilePath)
            .Select(i => i.Split(' ')
                .Select(int.Parse)
                .ToArray())
            .ToArray();
        
        PredictNextVal(_sequences[0]);
    }

    public override ValueTask<string> Solve_1()
        => new(_sequences.Select(PredictNextVal).Sum().ToString());

    public override ValueTask<string> Solve_2()
        => new(_sequences.Select(i => i.Reverse()).Select(PredictNextVal).Sum().ToString());

    private int PredictNextVal(IEnumerable<int> sequence)
    {
        var curList = sequence.ToArray();
        var lasts = new List<int> { curList[^1] };
        while (true)
        {
            curList = curList.Pairwise().Select(t => t.Item2 - t.Item1).ToArray();
            lasts.Add(curList[^1]);
            if (curList.All(i => i == 0))
                break;
        }
        
        for (int i = lasts.Count - 2; i >= 0; i--)
            lasts[i] += lasts[i + 1];

        return lasts[0];
    }
}
