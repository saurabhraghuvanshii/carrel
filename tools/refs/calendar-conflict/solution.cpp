#include <vector>
using namespace std;

int clashes(vector<vector<int>>& meetings, int start, int end) {
    int count = 0;
    for (auto& m : meetings)
        if (m[0] < end && start < m[1]) count++;
    return count;
}
