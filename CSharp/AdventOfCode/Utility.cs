﻿using System.Runtime;

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
    public static IEnumerable<(T, T)> Pairwise<T>(this IEnumerable<T> source)
    {
        var previous = default(T);
        using var it = source.GetEnumerator();
        if (it.MoveNext())
            previous = it.Current;

        while (it.MoveNext())
            yield return (previous, previous = it.Current);
    }
    // public static IEnumerable<(T val1, T val2)> Window<T>(this IReadOnlyList<T> enumerable)
    //     => enumerable.Zip(enumerable.Skip(1), (a, b) => (a, b));
}