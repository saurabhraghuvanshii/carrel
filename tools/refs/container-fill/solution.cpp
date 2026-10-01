#include <algorithm>
#include <vector>
using namespace std;

int mostWater(vector<int>& heights) {
    int best = 0, i = 0, j = (int)heights.size() - 1;
    while (i < j) {
        best = max(best, min(heights[i], heights[j]) * (j - i));
        if (heights[i] < heights[j]) i++;
        else j--;
    }
    return best;
}
