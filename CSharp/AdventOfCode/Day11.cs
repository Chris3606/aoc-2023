using SadRogue.Primitives;
using SadRogue.Primitives.GridViews;

namespace AdventOfCode;

public sealed class Day11 : BaseDay
{
    private readonly ArrayView<char> _grid;

    public Day11()
    {
        _grid = ParseUtils.CharGridFrom(File.ReadAllText(InputFilePath));
    }

    public override ValueTask<string> Solve_1()
        => new(
            ExpandGalaxy(_grid, 2)
                .ToArray()
                .Combinate()
                .Select(pair => Math.Abs(pair.Item2.X - pair.Item1.X) + Math.Abs(pair.Item2.Y - pair.Item1.Y))
                .Sum()
                .ToString());

    public override ValueTask<string> Solve_2()
        => new(
            ExpandGalaxy(_grid, 1000000)
                .ToArray()
                .Combinate()
                .Select(pair => Math.Abs(pair.Item2.X - pair.Item1.X) + Math.Abs(pair.Item2.Y - pair.Item1.Y))
                .Sum()
                .ToString());

    private static IEnumerable<Point64> ExpandGalaxy(IGridView<char> galaxy, int scaleFactor)
    {
        var xMap = new List<int>();
        var yMap = new List<int>();

        int curVal = 0;
        for (int y = 0; y < galaxy.Height; y++)
        {
            yMap.Add(curVal);
            
            if (Enumerable.Range(0, galaxy.Width).All(x => galaxy[x, y] == '.'))
                curVal += scaleFactor - 1;
            
        }

        curVal = 0;
        for (int x = 0; x < galaxy.Width; x++)
        {
            xMap.Add(curVal);
            
            if (Enumerable.Range(0, galaxy.Height).All(y => galaxy[x, y] == '.'))
                curVal += scaleFactor - 1;
        }
        
        return galaxy.Positions().Where(p => galaxy[p] == '#').Select(p => new Point64((long)p.X + xMap[p.X], (long)p.Y + yMap[p.Y]));
    }
}
