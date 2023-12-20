using System.Text;
using SadRogue.Primitives;
using SadRogue.Primitives.GridViews;

namespace AdventOfCode;

public sealed class Day14: BaseDay
{
    private readonly ArrayView<char> _grid;

    public Day14()
    {
        _grid = ParseUtils.CharGridFrom(File.ReadAllText(InputFilePath));
    }

    public override ValueTask<string> Solve_1()
    {
        SlideRocks(_grid,_grid.PositionsRowAscending(), Direction.Up);
        return new(_grid.Positions().Where(i => _grid[i] == 'O').Sum(i => _grid.Height - i.Y).ToString());
    }

    public override ValueTask<string> Solve_2()
    {
        const int cycles = 1000000000;
        
        Dictionary<string, int> seen = new(){ [SerializeRocks(_grid)] = 0};
        int period = 0;
        int curCycle = 0;
        for (int i = 1; i <= cycles; i++)
        {
            DoCycle(_grid);
            var serialized = SerializeRocks(_grid);
            if (seen.TryGetValue(serialized, out var prevCycle))
            {
                period = i - prevCycle;
                curCycle = i;
                break;
            }
            
            seen.Add(serialized, i);
        }
        
        if (period != 0)
        {
            int remaining = cycles - curCycle;
            remaining %= period;
            
            for (int i = 0; i < remaining; i++)
                DoCycle(_grid);
        }
        
        return new(_grid.Positions().Where(i => _grid[i] == 'O').Sum(i => _grid.Height - i.Y).ToString());
    }

    private static void SlideRocks(ISettableGridView<char> grid, IEnumerable<Point> points, Direction direction)
    {
        foreach (var pos in points)
        {
            if (grid[pos] != 'O')
                continue;
            
            var newPos = pos + direction;
            while (grid.Contains(newPos) && grid[newPos] == '.')
                newPos += direction;

            grid[pos] = '.';
            grid[newPos - direction] = 'O';
        }
    }

    private static void DoCycle(ISettableGridView<char> grid)
    {
        SlideRocks(grid, grid.PositionsRowAscending(), Direction.Up);
        SlideRocks(grid, grid.PositionsColAscending(), Direction.Left);
        SlideRocks(grid, grid.PositionsRowDescending(), Direction.Down);
        SlideRocks(grid,grid.PositionsColDescending(), Direction.Right);
    }
    
    private static string SerializeRocks(IGridView<char> grid)
    {
        var sb = new StringBuilder();
        foreach (var pos in grid.Positions().Where(i => grid[i] == 'O'))
            sb.Append(pos.X).Append(',').Append(pos.Y).Append(' ');
        return sb.ToString();
    }
}
