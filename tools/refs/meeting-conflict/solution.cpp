#include <algorithm>
#include <vector>
using namespace std;

bool canAttendAll(vector<vector<int>>& meetings) {
    vector<vector<int>> m = meetings;
    sort(m.begin(), m.end());
    for (size_t i = 1; i < m.size(); i++)
        if (m[i][0] < m[i - 1][1]) return false;
    return true;
}
