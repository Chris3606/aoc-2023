using System.Runtime;

namespace AdventOfCode;

public static class Utility
{
    public static Dictionary<T, int> BuildHistogram<T>(this IEnumerable<T> self)
    {
        var dict = new Dictionary<T, int>();
        foreach (var item in self)
        {
            if (!dict.TryAdd(item, 1))
                dict[item]++;
        }

        return dict;
    }

    // Gets greatest common divisor (GCD) via Euclidean algorithm
    public static long GCD(long a, long b)
    {
        while (b != 0)
        {
            long t = b;
            b = a % b;
            a = t;
        }

        return a;
    }

    public static long LCM(long a, long b, params long[] integers)
    {
        long result = a * b / GCD(a, b);

        foreach (var i in integers)
            result = LCM(result, i);

        return result;
    }
}