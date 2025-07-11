# Player Access Restrictions

## Updated Permission System

Players now have **very limited access** and can only manage match-related data.

### What Players CAN Do:
- ✅ **View their profile**: `GET /api/profile`
- ✅ **View public data**: provinces, parlours, posts (read-only)
- ✅ **Match operations**:
  - View matches (`GET /api/matches`)
  - View specific match (`GET /api/matches/:id`)
  - Create new matches (`POST /api/matches`)
  - Update match details (`PUT /api/matches/:id`)
- ✅ **Register** (public): `POST /api/players`

### What Players CANNOT Do:
- ❌ **View other players' data**: `GET /api/players`, `GET /api/players/:id`
- ❌ **Update player profiles**: `PUT /api/players/:id` (admin only)
- ❌ **Create/update provinces**: `POST /api/provinces`, `PUT /api/provinces/:id`
- ❌ **Create/update parlours**: `POST /api/parlours`, `PUT /api/parlours/:id`
- ❌ **Create/update posts**: `POST /api/posts`, `PUT /api/posts/:id`
- ❌ **Delete matches**: `DELETE /api/matches/:id` (admin only)
- ❌ **Approve matches**: `POST /api/matches/:id/approve` (admin only)
- ❌ **Access admin/role management**: All admin and role endpoints

## Key Changes Made:

1. **Removed player access** from all player management endpoints
2. **Removed player access** from all province/parlour/post modification endpoints
3. **Kept player access** only for match-related operations
4. **Players can only read public data** (provinces, parlours, posts)

## Security Benefits:

- **Data isolation**: Players cannot access other players' information
- **Limited modification rights**: Players can only modify match data
- **No administrative access**: Players cannot perform admin operations
- **Focused functionality**: Players can only do what they need for matches

## For Development:

- All changes are backward compatible
- No database changes required
- Only route-level access control updated
- Documentation updated to reflect new permissions

This ensures players can only manage match data and cannot post or update any other data tables in the system.
