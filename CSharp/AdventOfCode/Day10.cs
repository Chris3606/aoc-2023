using SadRogue.Primitives;
using SadRogue.Primitives.GridViews;

namespace AdventOfCode;

public sealed class Day10 : BaseDay
{
    private static readonly Dictionary<char, Direction[]> PipeNeighborMap = new()
    {
        ['|'] = new[] { Direction.Up, Direction.Down },
        ['-'] = new[] { Direction.Right, Direction.Left },
        ['L'] = new[] { Direction.Up, Direction.Right },
        ['J'] = new[] { Direction.Up, Direction.Left },
        ['7'] = new[] { Direction.Down, Direction.Left },
        ['F'] = new[] { Direction.Down, Direction.Right }
    };
    private readonly ArrayView<char> _grid;
    private readonly Point _start;

    public Day10()
    {
        _grid = ParseUtils.CharGridFrom(File.ReadAllText(InputFilePath));
        _start = _grid.Positions().First(i => _grid[i] == 'S');

        var startNeighbors = AdjacencyRule.Cardinals.DirectionsOfNeighborsClockwiseCache
            .Where(d => _grid.Contains(_start + d) && _grid[_start + d] != '.')
            .Where(d => PipeNeighborMap[_grid[_start + d]].Contains(d + 4));
        
        var startingPipe = PipeNeighborMap.First(pair => pair.Value.All(startNeighbors.Contains)).Key;
        _grid[_start] = startingPipe;
    }

    public override ValueTask<string> Solve_1()
        => new(GetMainLoop(_grid, _start).Values.Max().ToString());

    public override ValueTask<string> Solve_2()
    {
        var distances = GetMainLoop(_grid, _start);
        var bounds = distances.Keys.GetBounds();

        int innerPoints = 0;
        for (int y = bounds.MinExtentY; y <= bounds.MaxExtentY; y++)
        {
            int crosses = 0;
            for (int x = bounds.MinExtentX; x <= bounds.MaxExtentX; x++)
            {
                var pos = new Point(x, y);
                char c = distances.ContainsKey(pos) ? _grid[pos] : '.';
                
                switch (c)
                {
                    case 'F':
                        crosses++;
                        break;
                    case '7':
                        crosses++;
                        break;
                    case '|':
                        crosses++;
                        break;
                    case '-':
                        break;
                    case '.':
                        if (crosses % 2 == 1)
                            innerPoints++;
                        break;
                }
            }
        }

        return new(innerPoints.ToString());
    }

    private static Dictionary<Point, int> GetMainLoop(IGridView<char> grid, Point start)
    {
        var q = new Queue<(Point pos, int dist)>();
        q.Enqueue((start, 0));
        var distances = new Dictionary<Point, int>();
        
        while (q.Count > 0)
        {
            var cur = q.Dequeue();
            if (!distances.TryAdd(cur.pos, cur.dist))
                continue;

            foreach (var dir in PipeNeighborMap[grid[cur.pos]])
            {
                var next = cur.pos + dir;
                if (grid.Contains(next) && grid[next] != '.')
                    q.Enqueue((next, cur.dist + 1));
            }
        }

        return distances;
    }
}
